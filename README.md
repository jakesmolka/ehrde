# EHRDE - (open) Electronic Health Record Data Explorer

## Disclaimer

This repository is only a historical documentation of my bachelor thesis's code. It's non-maintained research code.

Additionally, this code is about one year old on the day of this publication. It may not represent how I would solve specific problems as of today.

### About
A Dashboard to explore openEHR data hosted within a Think!EHR Platform instance.

### Requirements
 - Go, installed Go(lang) environment
 - Insert jqwidget files into assets/js/jqwidgets

### Installation
 1. Extract archive (source files and assets)
 2. Run ``go run main.go``

### Config
Edit ``config.json`` file according to your:
 - BaseUrl: Base url to access Think!EHR server's REST service including trailing ``/``
 - User: Username for HTTP basic authentication
 - Pass: Password for HTTP basic authentication
