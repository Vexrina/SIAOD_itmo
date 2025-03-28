package internal

import (
	"context"
	desc "vexrina/siaod_itmo/lab_03/lab_03/pkg/lab_03"
)

type (
	Usecase interface {
		LoadFile(req *desc.LoadFileRequest) (*desc.LoadFileReply, error)
		GetAllFiles(req *desc.GetAllFilesRequest) (*desc.GetAllFilesResponse, error)
		FindTerm(req *desc.FindTermRequest) (*desc.FindTermResponse, error)
	}

	service struct {
		desc.UnimplementedLab03Server

		usecase Usecase
	}
)

func NewService(usecase Usecase) service {
	return service{usecase: usecase}
}

func (s service) LoadFile(ctx context.Context, req *desc.LoadFileRequest) (*desc.LoadFileReply, error) {
	_ = ctx
	return s.usecase.LoadFile(req)
}

func (s service) FindTerm(ctx context.Context, req *desc.FindTermRequest) (*desc.FindTermResponse, error) {
	_ = ctx
	return s.usecase.FindTerm(req)
}

func (s service) GetAllFiles(ctx context.Context, req *desc.GetAllFilesRequest) (*desc.GetAllFilesResponse, error) {
	_ = ctx
	return s.usecase.GetAllFiles(req)
}
