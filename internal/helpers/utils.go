package helpers

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"entain/internal/domain"
	errors "entain/internal/error"

	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
)

// UsageFor - print usage.
func UsageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

// LoadConfig - reads configuration from the environment variables.
func LoadConfig() (config domain.AppConfig, _ error) {
	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("failed to read configuration: %v", err)
	}
	return config, nil
}

// ValidateConfig - validates configuration.
func ValidateConfig(config domain.AppConfig) error {

	err := validator.New().Struct(config)
	if err != nil {

		// This check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return errors.NewErrFailedPrecondition(
				fmt.Sprintf("%s required", err.Field()),
			)
		}
	}
	return nil
}
