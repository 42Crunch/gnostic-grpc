// Copyright 2021 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is an aggregation tool which scans for and aggregates incompatibilities
// in OpenAPI documents within a given directory. The only argument given to this
// tool should be the intended directory.

package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/googleapis/gnostic-grpc/incompatibility"
	"github.com/googleapis/gnostic-grpc/utils"
)

// main function for aggreation tool
func main() {
	if len(os.Args) != 2 {
		exitIfError(errors.New("argument should be a path to a directory"))
	}
	generateAnalysis(os.Args[1])
}

// runs analysis on given directory
func generateAnalysis(dirPath string) *incompatibility.ApiSetIncompatibility {
	analysisAggregation := incompatibility.NewAnalysis()
	readingDirectoryErr := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("walk error for file at %s", path)
		}
		newAnalysis, analysisErr := router(dirPath, path, d)
		if analysisErr != nil {
			log.Printf("unable to produce analysis for file %s with error <%s>", path, analysisErr.Error())
		} else {
			analysisAggregation = incompatibility.AggregateAnalysis(analysisAggregation, newAnalysis)
		}
		return nil
	})
	if readingDirectoryErr != nil {
		exitIfError(errors.New("unable to walk through directory"))
	}
	return analysisAggregation
}

// router directs logic for either a file or a directory to produce an analysis in either case
func router(parentDir string, path string, d fs.DirEntry) (*incompatibility.ApiSetIncompatibility, error) {
	if path == parentDir {
		return incompatibility.NewAnalysis(), nil
	}
	if d.IsDir() {
		return directoryHandler(path)
	}
	singleFileReport, reportErr := fileHandler(path)
	if reportErr != nil {
		return nil, reportErr
	}
	return incompatibility.AggregateReports(singleFileReport), nil
}

// fileHander attempts to parse the file at path and to then create an incompatibility report
func fileHandler(path string) (*incompatibility.IncompatibilityReport, error) {
	openAPIDoc, err := utils.ParseOpenAPIDoc(path)
	if err != nil {
		return nil, err
	}
	incompatibilityReport := incompatibility.ScanIncompatibilities(openAPIDoc)
	log.Printf("created incompatibility report for file at %s\n", path)
	return incompatibilityReport, nil
}

// TODO
// directoryHandler attempts to generate an ApiSetIncompatibilty object containing
// incompatibility information of all openapi files in the given directory
func directoryHandler(dirPath string) (*incompatibility.ApiSetIncompatibility, error) {
	return nil, errors.New("unable to currently handle directories")
}

func exitIfError(e error) {
	if e == nil {
		return
	}
	log.Printf("Exiting with error: %s\n", e)
	os.Exit(1)
}
