/*
Copyright 2023 The Cloud-Barista Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package s3fs

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/cloud-barista/mc-data-manager/pkg/utils"
)

type reader struct {
	r        *io.PipeReader
	ch       chan error
	cancel   context.CancelFunc
	chkClose bool
}

func (p *reader) Read(b []byte) (int, error) {
	return p.r.Read(b)
}

func (p *reader) Close() error {
	if !p.chkClose {
		p.chkClose = true
		return p.r.Close()
	}
	return nil
}

type writer struct {
	w        *io.PipeWriter
	ch       chan error
	cancel   context.CancelFunc
	chkClose bool
}

func (p *writer) Write(b []byte) (int, error) {
	return p.w.Write(b)
}

func (p *writer) Close() error {
	if !p.chkClose {
		p.chkClose = true
		_ = p.w.Close()
		return <-p.ch
	}
	return nil
}

type fakeWriteAt struct {
	W io.Writer
}

func (w *fakeWriteAt) WriteAt(p []byte, off int64) (n int, err error) {
	return w.W.Write(p)
}

type S3FS struct {
	provider   utils.Provider
	bucketName string
	region     string

	client     *s3.Client
	ctx        context.Context
	uploader   manager.Uploader
	downloader manager.Downloader
}

// Creating a Bucket
//
// Aws imposes location constraints when creating buckets
func (f *S3FS) CreateBucket() error {
	_, err := f.client.HeadBucket(f.ctx, &s3.HeadBucketInput{
		Bucket: aws.String(f.bucketName),
	})

	if err != nil {
		var bae *types.BucketAlreadyExists
		if errors.As(err, &bae) {
			return nil
		}
		var baoby *types.BucketAlreadyOwnedByYou
		if errors.As(err, &baoby) {
			return nil
		}
		var nf *types.NotFound
		var nsb *types.NoSuchBucket
		if errors.As(err, &nsb) || errors.As(err, &nf) {
			input := &s3.CreateBucketInput{Bucket: aws.String(f.bucketName)}
			if f.provider == "aws" {
				input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
					LocationConstraint: types.BucketLocationConstraint(f.region),
				}
			}
			_, err := f.client.CreateBucket(f.ctx, input)
			return err
		}
		return err
	}
	return err
}

// Delete Bucket
//
// Check and delete all objects in the bucket and delete the bucket
func (f *S3FS) DeleteBucket() error {
	objList, err := f.ObjectList()
	if err != nil {
		return err
	}

	if len(objList) != 0 {
		var objectIds []types.ObjectIdentifier
		for _, object := range objList {
			objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(object.Key)})
		}

		_, err = f.client.DeleteObjects(f.ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(f.bucketName),
			Delete: &types.Delete{Objects: objectIds},
		})

		if err != nil {
			return err
		}
	}
	_, err = f.client.DeleteBucket(f.ctx, &s3.DeleteBucketInput{Bucket: &f.bucketName})
	if err != nil {
		return err
	}
	return nil
}

// Open function using pipeline
func (f *S3FS) Open(name string) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	ch := make(chan error)
	ctx, cancel := context.WithCancel(f.ctx)
	go func() {
		defer cancel()
		_, err := f.downloader.Download(
			ctx,
			&fakeWriteAt{W: pw},
			&s3.GetObjectInput{
				Bucket: aws.String(f.bucketName),
				Key:    aws.String(name),
			}, func(d *manager.Downloader) { d.Concurrency = 1 },
		)
		if cerr := pw.Close(); cerr != nil {
			err = cerr
		}
		ch <- err
	}()

	return &reader{r: pr, ch: ch, cancel: cancel, chkClose: false}, nil
}

// Create function using pipeline
func (f *S3FS) Create(name string) (io.WriteCloser, error) {
	pr, pw := io.Pipe()
	ch := make(chan error)
	ctx, cancel := context.WithCancel(f.ctx)
	go func() {
		defer cancel()
		_, err := f.uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(f.bucketName),
			Key:    aws.String(name),
			Body:   pr,
		})
		ch <- err
	}()

	return &writer{w: pw, ch: ch, cancel: cancel, chkClose: false}, nil
}

// Look up the list of objects in your bucket
func (f *S3FS) ObjectList() ([]*utils.Object, error) {
	var objlist []*utils.Object
	var ContinuationToken *string

	for {
		LOut, err := f.client.ListObjectsV2(
			f.ctx,
			&s3.ListObjectsV2Input{
				Bucket:            aws.String(f.bucketName),
				ContinuationToken: ContinuationToken,
			},
		)
		if err != nil {
			return nil, err
		}

		for _, obj := range LOut.Contents {
			objlist = append(objlist, &utils.Object{
				ETag:         *obj.ETag,
				Key:          *obj.Key,
				LastModified: *obj.LastModified,
				Size:         *obj.Size,
				StorageClass: string(obj.StorageClass),
			})
		}

		if LOut.NextContinuationToken == nil {
			break
		}

		ContinuationToken = LOut.NextContinuationToken
	}

	return objlist, nil
}

func New(provider utils.Provider, client *s3.Client, bucketName, region string) *S3FS {
	sfs := &S3FS{
		ctx:        context.TODO(),
		provider:   provider,
		bucketName: bucketName,
		region:     region,
		client:     client,
	}

	sfs.uploader = *manager.NewUploader(client, func(u *manager.Uploader) { u.Concurrency = 1; u.PartSize = 128 * 1024 * 1024 })
	sfs.downloader = *manager.NewDownloader(client, func(d *manager.Downloader) { d.Concurrency = 1; d.PartSize = 128 * 1024 * 1024 })

	return sfs
}
