# Release Notes
Another day, another release. These are the release notes for the version `v0.22.2`.

This release introduces a new argument in both vulnerability exception resources, inside `vulnerability_criteria`
named `fixable_vuln`, this new argument is of type `string` and it solves the tristate boolean issue reported at
https://github.com/lacework/terraform-provider-lacework/issues/339.

The (tres) states are:
* `"false"` (Yes! A string with value `false`) - Constraints the exception by **non-fixable** vulnerabilities
* `"true"` (A string with value `true`) - Constraints the exception by **fixable** vulnerabilities
* `""` (empty string) - **Do not filter** by fixable or non-fixable vulnerability

**Note** that we deprecated the previous boolean argument `fixable` and it will be removed on version `v1.0`.

## Other Changes
* build(deps): bump github.com/stretchr/testify from 1.7.3 to 1.8.0 (#337) (dependabot[bot])([2229d06](https://github.com/lacework/terraform-provider-lacework/commit/2229d069c887086f4a9e5a4f23fef6ce8b14b586))
