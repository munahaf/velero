/*
Copyright 2017 the Velero contributors.

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

package v2

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"

	api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
)

// BackupItemAction is an actor that performs an operation on an individual item being backed up.
type BackupItemAction interface {
	// AppliesTo returns information about which resources this action should be invoked for.
	// A BackupItemAction's Execute function will only be invoked on items that match the returned
	// selector. A zero-valued ResourceSelector matches all resources.
	AppliesTo() (velero.ResourceSelector, error)

	// Execute allows the BackupItemAction to perform arbitrary logic with the item being backed up,
	// including mutating the item itself prior to backup. The item (unmodified or modified)
	// should be returned, along with an optional slice of ResourceIdentifiers specifying
	// additional related items that should be backed up.
	Execute(item runtime.Unstructured, backup *api.Backup) (runtime.Unstructured, []velero.ResourceIdentifier, string, error)

	// Progress allows the BackupItemAction to report on progress of an asynchronous action.
	// For the passed-in operation, the plugin will return an OperationProgress struct, indicating
	// whether the operation has completed, whether there were any errors, a plugin-specific
	// indication of how much of the operation is done (items completed out of items-to-complete),
	// and started/updated timestamps
	Progress(operationID string, backup *api.Backup) (velero.OperationProgress, error)

	// Cancel allows the BackupItemAction to cancel an asynchronous action (if possible).
	// Velero will call this if the wait timeout for asynchronous actions has been reached.
	// If operation cancel is not supported, then the plugin just needs to return. No error
	// return is expected in this case, since cancellation is optional here.
	Cancel(operationID string, backup *api.Backup) error
}

func AsyncOperationsNotSupportedError() error {
	return errors.New("Plugin does not support asynchronous operations")
}

func InvalidOperationIDError(operationID string) error {
	return errors.New(fmt.Sprintf("Operation ID %v is invalid.", operationID))
}
