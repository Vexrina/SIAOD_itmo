package api

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"log"
	"strings"
	"vexrina/siaod_itmo/lab_03/internal/utils"
	desc "vexrina/siaod_itmo/lab_03/lab_03/pkg/lab_03"
)

func (l lab03usecase) FindTerm(req *desc.FindTermRequest) (*desc.FindTermResponse, error) {
	index, err := bleve.Open("bleve/" + req.FileName + ".bleve")
	if err != nil {
		log.Println(err)
		return &desc.FindTermResponse{}, err
	}
	defer index.Close()
	res, err := utils.SearchIndex(index, strings.Split(req.Term, " ")...)
	if err != nil {
		log.Println(err)
		return &desc.FindTermResponse{}, err
	}
	result := fmt.Sprintf("%v\n", res.Total)
	for _, hit := range res.Hits {
		result = result + fmt.Sprintf("ID: %s, Score: %.2f\n", hit.ID, hit.Score)
	}
	return &desc.FindTermResponse{Message: result}, nil
}
