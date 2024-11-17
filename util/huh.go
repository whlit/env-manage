package util

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/logger"
)

func Select[T comparable](getKey func(T) string, items ...T) (T, bool) {
	var selected T
	var options []huh.Option[T]
	for _, item := range items {
		options = append(options, huh.NewOption(getKey(item), item))
	}
	err := huh.NewSelect[T]().Options(options...).Value(&selected).Run()
	if err != nil {
		logger.Info("选择失败", err)
		return selected, false
	}
	return selected, true
}

func SelectWithConfirm[T comparable](getKey func(T) string, items ...T) (T, bool) {
	var selected T
	var options []huh.Option[T]
	for _, item := range items {
		options = append(options, huh.NewOption(getKey(item), item))
	}
	err := huh.NewSelect[T]().Options(options...).Value(&selected).Run()
	if err != nil {
		logger.Info("选择失败", err)
		return selected, false
	}
	var confirm bool
	huh.NewConfirm().Title(strings.Join([]string{"确认选择 ", getKey(selected), " ?"}, "")).Value(&confirm).Run()
	return selected, confirm
}