package api

import (
	"strings"
	"vexrina/siaod_itmo/lab_03/internal/utils"
	desc "vexrina/siaod_itmo/lab_03/lab_03/pkg/lab_03"
)

func (l lab03usecase) GetAllFiles(_ *desc.GetAllFilesRequest) (*desc.GetAllFilesResponse, error) {
	index, not_index, err := utils.FindMatchingFiles("csv/", "bleve/")
	if err != nil {
		return &desc.GetAllFilesResponse{}, nil
	}
	return &desc.GetAllFilesResponse{Indexed: strings.Join(index, ","), NotIndexed: strings.Join(not_index, ",")}, nil
}
