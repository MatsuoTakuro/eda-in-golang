package customers

//go:generate buf generate

//go:generate mockery --quiet --dir ./customerspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore
