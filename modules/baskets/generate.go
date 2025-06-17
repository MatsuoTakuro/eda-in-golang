package baskets

//go:generate buf generate

//go:generate mockery --quiet --dir ./basketspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore
