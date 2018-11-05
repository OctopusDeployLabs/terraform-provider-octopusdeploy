# Decisions

## Development

* You cannot use `ConflictsWith` inside `schema.TypeList`. Instead I am splitting them up into their own `schema.TypeList`.

## Documentation

* Documentation for usage of the Terraform provider will go into the `/docs/provider/` folder and be split into the `data_sources` or `resources` sub-folder.

* Documentation will be done in Markdown and match the style of the official Terraform providers. This is to make it easy to transition to official provider documentation if this occurs.
