package ordering

//go:generate buf generate

//go:generate mockery --quiet --dir ./orderingpb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore
