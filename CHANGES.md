# syngo - Changes <!-- omit in toc -->


## 0.3.1 - 20th August 2026

* enforced Synesis Go import order via **gci** (**.golangci.yml**, **examples/.golangci.yml**);
* normalised example companion docs to the per-program layout (`examples/<name>/main.go` + `examples/<name>.md`);
* version string updated for the 0.3.1 release;


## 0.3.0 - 20th August 2026

* added **Version()** (replacing the **Version** constant), formed by **ver2go.CombineVersion()**;
* updated **ver2go** to 0.2.0-beta1;
* version string updated for the 0.3.0 release;


## 0.2.1 - 20th August 2026

* CI modernisation (matrix + lint);
* boilerplate additions (scripts, markdown docs, project identity);
* version string updated for the 0.2.1 release;


## 0.2.0 - 4th August 2026

* added `DownCounter` and `UpCounter` types;
* added GitHub Actions CI (**go.yml**);
* updated **build.sh** and **run_all_unit_tests.sh**;
* lowered Go toolchain requirement to 1.23.6;
* fixed latch tests for Go versions without `WaitGroup.Go`;
* boilerplate and documentation updates;


## 0.2.0-alpha1 - 2nd November 2025

FIRST PUBLIC RELEASE



<!-- ########################### end of file ########################### -->
