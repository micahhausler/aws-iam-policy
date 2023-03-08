COVERAGE=cover.out

$(COVERAGE):
	go test -v -coverprofile=$(COVERAGE) ./...

test: $(COVERAGE)

cover: test
	go tool cover -func=$(COVERAGE)

clean:
	go clean -testcache
	rm -rf ./$(COVERAGE)


# This downloads all the latest AWS IAM policies and stores them in the
# managed-policies directory. We don't store them in the repo because they
# change frequently change, take up ~4.8MB, and include 1050 different files.
update-fixtures:
	mkdir -p policy/test_fixtures/managed-policies/
	for POLICY in $$(aws iam list-policies --scope AWS | jq -c '.Policies[] | select(.IsAttachable==true) | {version: .DefaultVersionId, name: .PolicyName, arn: .Arn}'); do \
	 echo $$POLICY | jq .name; \
	 aws iam get-policy-version \
		--policy-arn $$(echo $$POLICY | jq -r .arn) \
		--version-id $$(echo $$POLICY | jq  -r .version) | \
			jq -Sr '.PolicyVersion.Document' > policy/test_fixtures/managed-policies/$$(echo $$POLICY | jq -r .name).json; \
	done

.PHONY: test cover clean update-fixtures
