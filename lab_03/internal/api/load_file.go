package api

import (
	"log"

	"github.com/blevesearch/bleve/v2"

	"vexrina/siaod_itmo/lab_03/internal/utils"
	desc "vexrina/siaod_itmo/lab_03/lab_03/pkg/lab_03"
)

func (l lab03usecase) LoadFile(req *desc.LoadFileRequest) (*desc.LoadFileReply, error) {
	csvName := req.Name

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("bleve/"+csvName+".bleve", mapping)
	if err != nil {
		index, err = bleve.Open("bleve/" + csvName + ".bleve")
		if err != nil {
			log.Fatal(err)
		}
		defer func(index bleve.Index) {
			err = index.Close()
			if err != nil {

			}
		}(index)
	}

	err = utils.ReadCsv(req.Name, index)
	if err != nil {
		return &desc.LoadFileReply{}, err
	}

	return &desc.LoadFileReply{}, nil
}
