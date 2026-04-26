//go:build cublas

package llama

/*
#cgo CPPFLAGS: -DGGML_USE_CUDA
#cgo LDFLAGS: -lggml-cuda -lcublas -lcudart -L/usr/local/cuda/lib64/
*/
import "C"
