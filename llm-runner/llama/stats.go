package llama

/*
#include "wrapper.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

type GPUInfo struct {
	DeviceID      int    // id CUDA-устройства
	DeviceName    string // Название модели gpu (например, "Quadro RTX 6000")
	FreeMemoryMB  int    // Доступная VRAM в МБ
	TotalMemoryMB int    // Общий объем VRAM в МБ
}

type ModelMetadata struct {
	Architecture string // Архитектура модели (например, qwen3, llama)
	Name         string // Полное имя модели
	Basename     string // Базовое имя модели
	QuantizedBy  string // Кем выполнена квантизация модели
	SizeLabel    string // Размер модели (например, 8B, 70B)
	RepoURL      string // Адрес репозитория на Hugging Face
}

type RuntimeInfo struct {
	ContextSize     int    // Размер контекстного окна в токенах
	BatchSize       int    // Размер пакетной обработки
	KVCacheType     string // Тип квантизации KV-кэша ("f16", "q8_0", "q4_0")
	KVCacheSizeMB   int    // Оценка использования памяти KV-кэша в МБ
	GPULayersLoaded int    // Количество слоев, выгруженных на GPU
	TotalLayers     int    // Общее количество слоев в модели
}

type ModelStats struct {
	GPUs     []GPUInfo     // Информация о доступных CUDA GPU
	Metadata ModelMetadata // Метаданные модели из GGUF-файла
	Runtime  RuntimeInfo   // Конфигурация рантайма и использование ресурсов
}

func (m *Model) Stats() (*ModelStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, fmt.Errorf("модель закрыта")
	}

	stats := &ModelStats{}

	gpuCount := int(C.llama_wrapper_get_gpu_count())
	stats.GPUs = make([]GPUInfo, 0, gpuCount)

	for i := 0; i < gpuCount; i++ {
		var cInfo C.llama_wrapper_gpu_info
		if C.llama_wrapper_get_gpu_info(C.int(i), &cInfo) {
			stats.GPUs = append(stats.GPUs, GPUInfo{
				DeviceID:      int(cInfo.device_id),
				DeviceName:    C.GoString(&cInfo.device_name[0]),
				FreeMemoryMB:  int(cInfo.free_memory_mb),
				TotalMemoryMB: int(cInfo.total_memory_mb),
			})
		}
	}

	stats.Metadata = ModelMetadata{
		Architecture: m.getMetaString("general.architecture"),
		Name:         m.getMetaString("general.name"),
		Basename:     m.getMetaString("general.basename"),
		QuantizedBy:  m.getMetaString("general.quantized_by"),
		SizeLabel:    m.getMetaString("general.size_label"),
		RepoURL:      m.getMetaString("general.repo_url"),
	}

	return stats, nil
}

func (m *Model) getMetaString(key string) string {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cValue := C.llama_wrapper_model_meta_string(m.modelPtr, cKey)
	if cValue == nil {
		return ""
	}

	return C.GoString(cValue)
}

func (s *ModelStats) String() string {
	var b strings.Builder

	b.WriteString("=== Model Statistics ===\n\n")

	if len(s.GPUs) > 0 {
		b.WriteString("GPU Devices:\n")
		for _, gpu := range s.GPUs {
			fmt.Fprintf(&b, "  GPU %d: %s\n", gpu.DeviceID, gpu.DeviceName)
			fmt.Fprintf(&b, "    VRAM: %d MB free / %d MB total\n", gpu.FreeMemoryMB, gpu.TotalMemoryMB)
		}

		b.WriteString("\n")
	}

	b.WriteString("Model Details:\n")
	if s.Metadata.Name != "" {
		fmt.Fprintf(&b, "  Name: %s\n", s.Metadata.Name)
	}

	if s.Metadata.Architecture != "" {
		arch := s.Metadata.Architecture
		if s.Metadata.SizeLabel != "" {
			arch += " (" + s.Metadata.SizeLabel + ")"
		}

		fmt.Fprintf(&b, "  Architecture: %s\n", arch)
	}

	if s.Metadata.QuantizedBy != "" {
		fmt.Fprintf(&b, "  Quantized by: %s\n", s.Metadata.QuantizedBy)
	}

	if s.Metadata.RepoURL != "" {
		fmt.Fprintf(&b, "  Repository: %s\n", s.Metadata.RepoURL)
	}

	b.WriteString("\n")

	b.WriteString("Runtime Configuration:\n")
	fmt.Fprintf(&b, "  Context: %s tokens | Batch: %d tokens\n", formatNumber(s.Runtime.ContextSize), s.Runtime.BatchSize)
	fmt.Fprintf(&b, "  KV Cache: %s (%s MB)\n", s.Runtime.KVCacheType, formatNumber(s.Runtime.KVCacheSizeMB))
	fmt.Fprintf(&b, "  GPU Layers: %d/%d\n", s.Runtime.GPULayersLoaded, s.Runtime.TotalLayers)

	return b.String()
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	s := fmt.Sprintf("%d", n)
	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}

		result.WriteRune(c)
	}

	return result.String()
}
