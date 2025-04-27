name := "api-performance-tests"
version := "1.0"
scalaVersion := "2.13.10"

enablePlugins(GatlingPlugin)

libraryDependencies += "io.gatling.highcharts" % "gatling-charts-highcharts" % "3.9.5" % "test"
libraryDependencies += "io.gatling" % "gatling-test-framework" % "3.9.5" % "test"

// 在build.sbt中添加日志级别配置
javaOptions in Gatling += "-Dlogback.configurationFile=logback.xml"
