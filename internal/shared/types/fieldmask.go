package types

import "google.golang.org/protobuf/types/known/fieldmaskpb"

// FieldMaskContains checks if a field path is in the mask.
func FieldMaskContains(mask *fieldmaskpb.FieldMask, path string) bool {
	if mask == nil {
		return false
	}
	for _, p := range mask.Paths {
		if p == path {
			return true
		}
	}
	return false
}

// FieldMaskPaths returns all paths in the mask, or empty if nil.
func FieldMaskPaths(mask *fieldmaskpb.FieldMask) []string {
	if mask == nil {
		return nil
	}
	return mask.Paths
}

// NewFieldMask creates a new FieldMask with the given paths.
func NewFieldMask(paths ...string) *fieldmaskpb.FieldMask {
	return &fieldmaskpb.FieldMask{Paths: paths}
}
