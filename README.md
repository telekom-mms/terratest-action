# terraform-action

A GitHub action to run terratest test cases against your repository.

## Usage

### Prepare Repository

Create the following in your repository where the action should run.

First setup a `tests` directory and create the following files.

`tests/main.tf`

* terraform configuration to test with this action

`tests/terratest.yaml`

* configuration what tests should be run with terratest

### Setup Workflow

To use this GitHub Action, add the following parts to your YAML workflow file, e.g. `.github/workflows/terratest.yml`.

``` yaml
...

jobs:
  ...
    steps:
      ...

      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.19
          check-latest: true

      - name: Checkout Github Action
        uses: actions/checkout@v3
        with:
          repository: T-Systems-MMS/terratest-action
          ref: main
          path: ./terratest-action
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      ...

      - name: Run Test for ${{ matrix.example }} Setup
        uses: ./terratest-action
        with:
          test:  azure
          test_type: unit
          azure_client_id: ${{ secrets.AZURE_CLIENT_ID }}
          azure_client_secret: ${{ secrets.AZURE_CLIENT_SECRET }}
          azure_subscription_id: ${{ secrets.AZURE_SUBSCRIPTION_ID }}
          azure_tenant_id: ${{ secrets.AZURE_TENANT_ID }}
```

## Configuration

### Workflow

#### Inputs

| name                  | description                                                                         | required |
| --------------------- | ----------------------------------------------------------------------------------- | -------- |
| test                  | name of the test to run (e.g. azure)                                                | true     |
| test_type            | run unit and integration tests, set to `unit` if integration tests should be skipped | false    |
| azure_client_id       | application id of azure ad application                                              | false    |
| azure_client_secret   | application password of azure ad application                                        | false    |
| azure_subscription_id | azure subscription id                                                               | false    |
| azure_tenant_id       | azure tenant id                                                                     | false    |

### Terratest

| name      | description                                                                    | `terratest.yaml` |
| ----------| ------------------------------------------------------------------------------ | ---------------- |
| functions | functions that should be testet bases on the test (e.g. azure) that should run | `functions: []`  |
| options   | key: value pair that should be proved against the build                        | `options: {}`    |

#### azure

| test case | description                                                                       |
| ----------| --------------------------------------------------------------------------------- |
| .*Exists  | check if terraform resource exist, could be created                               |
| .*Show    | check if created terraform resource matches the conditions that should be checked |

Actually the following test cases are supported:

* ContainerRegistryExists
* ContainerRegistryShow

## Examples

* [Workflow configuration](examples/workflow.yml)
* [Terratest configuration](examples/terratest.yaml)
