package inference

import (
	"regexp"
	"strings"

	"github.com/sourcegraph/sourcegraph/lib/codeintel/autoindex/config"
)

func CanIndexJavaRepo(gitserver GitClient, paths []string) bool {
	return javaBuildTool(paths) != ""
}

func InferJavaIndexJobs(gitserver GitClient, paths []string) (indexes []config.IndexJob) {
	if buildTool := javaBuildTool(paths); buildTool != "" {
		indexes = append(indexes, config.IndexJob{
			Indexer: "sourcegraph/lsif-java",
			IndexerArgs: []string{
				"/coursier launch --contrib --ttl 0 lsif-java -- index --build-tool=" + buildTool,
			},
			Outfile: "dump.lsif",
			Root:    "",
			Steps:   []config.DockerStep{},
		})
	}
	return indexes
}

func JavaPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		suffixPattern(rawPattern("lsif-java.json")),
		suffixPattern(rawPattern(".java")),
		suffixPattern(rawPattern(".scala")),
	}
}

func javaBuildTool(paths []string) string {
	for _, buildToolPath := range paths {
		// The "lsif-java.json" file is generated by the JVMPACKAGES
		// external service type. This file is used to index package
		// repositories such as the JDK sources and published Java
		// libraries.
		if buildToolPath == "lsif-java.json" {
			for _, path := range paths {
				if isLsifJavaIndexablePath(path) {
					return "lsif"
				}
			}
			return ""
		}
		// Maven, Gradle and sbt are intentionally left out to begin
		// with as we gain more experience with auto-indexing package
		// repos, which have a higher likelyhood of indexing
		// successfully because they have a simpler build structure
		// compared to Gradle/Maven repos.
	}
	return ""
}

func isLsifJavaIndexablePath(path string) bool {
	return strings.HasSuffix(path, ".java") ||
		strings.HasSuffix(path, ".scala")
}
