/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	s3crdv1apiv1 "s3crdapiv1/api/v1"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3BucketReconciler reconciles a S3Bucket object
type S3BucketReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=s3crdv1api.s3crdapiv1,resources=s3buckets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=s3crdv1api.s3crdapiv1,resources=s3buckets/status,verbs=get;update;patch

func (r *S3BucketReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("s3bucket", req.NamespacedName)

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	// your logic here
	var s3bucket s3crdv1apiv1.S3Bucket

	if err := r.Get(ctx, req.NamespacedName, &s3bucket); err != nil {
		//log.Error(err, "unable to fetch EC2Instance")
		fmt.Println("unable to fetch s3bucket")
	}

	// Create the S3 Bucket
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s3bucket.Spec.BucketName),
	})
	if err != nil {
		fmt.Println(s3bucket.Spec.BucketName)
		fmt.Println("Unable to create bucket")
		fmt.Println(err)
	} else {
		fmt.Println("bucket creation success")
	}

	return ctrl.Result{}, nil
}

func (r *S3BucketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&s3crdv1apiv1.S3Bucket{}).
		Complete(r)
}
