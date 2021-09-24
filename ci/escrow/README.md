# Run

```shell script
brew install deno
deno run --allow-read --allow-write --allow-net ci/escrow/modify-pdf.ts 1.6.0-beta7 123456789 'Janice Morales' 'November 10, 2020' 'janice.morales@solo.io' '(617)-893-7557'
```

# Usage

Run the `generate-escrow-pdf` make target in order to generate an invoice for escrow deposit. The invoice will be stored at `_gloo-ee-source/ExhibitB<VERSION>.pdf`, where `<VERSION>` is the current semver. Additionally, a tar archive ove the repository's source will be created, and stored at `_gloo-ee-source/gloo-ee-$(VERSION).tar`

Relevant variables:
 |Variable Name|Default Value|Description|
 |-------------|-------------|-----------|
 |`TAGGED_VERSION`|The current Semver, as determined by the repo's `HEAD` commit|The version of Gloo Edge Enterprise that a deposit is being made for. Must follow [Semver format](https://semver.org/) and be prefixed with a literal "v", i.e., v1.8.1|
 |`DEPOSITOR_NAME`|Janice Morales|The value of the "Print Name" field in the "Deposit Certification" section of the invoice|
 |`DEPOSIT_DATE`|The current date, as determined by the output of the bash `date` command|The value of the "Date" field in the "Deposit Certification" section of the invoice|
 |`DEPOSITOR_EMAIL`|janice.morales@solo.io|The value of the "Email Address" field in the "Deposit Certification" section of the invoice|
 |`DEPOSITOR_PHONE`|(617)-893-7557|The value of the "Telephone Number" field in the "Deposit Certification" section of the invoice|

 Example Usage:
 ```sh
 $ TAGGED_VERSION=v1.8.1 \
 DEPOSITOR_NAME='Jane\ Doe' \
 DEPOSIT_DATE='Wed\ Dec\ 31\ 20:00:00\ EDT\ 1969' \
 DEPOSITOR_EMAIL='jane.doe@solo.io' \
 TELEPHONE_NUMBER='(555)-555-5555' \
 make generate-escrow-pdf
 ```

 Executing the above will create an invoice at `_gloo-ee-source/ExhibitB1.8.1.pdf` and an archive of the repository's source at `_gloo-ee-source/gloo-ee-1.8.1.tar`

 # Dependencies

 The `generate-escrow-pdf` make target uses the script located at [ci/escrow/modify-pdf.ts](ci/escrow/modify-pdf.ts) to populate data on the invoice. We use [Deno](https://deno.land/), a modern Javascript/Typescript runtime in order to execute this script. You can find install information for Deno at the following location: https://deno.land/#installation