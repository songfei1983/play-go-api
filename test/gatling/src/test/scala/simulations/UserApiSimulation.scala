package simulations

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._
import java.util.UUID

class UserApiSimulation extends Simulation {
  val kongHost = System.getenv.getOrDefault("KONG_HOST", "http://kong:8000")

  val httpProtocol = http
    .baseUrl(kongHost)
    .acceptHeader("application/json")
    .contentTypeHeader("application/json")

  val feeder = Iterator.continually(Map(
    "username" -> s"benchuser_${UUID.randomUUID().toString}",
    "password" -> "benchpass",
    "userId" -> "1"
  ))

  val registerAndLogin = feed(feeder)
    .exec(
      http("Register")
        .post("/api/v1/register")
        .body(StringBody("""{"username":"${username}","password":"${password}","email":"${username}@example.com"}"""))
        .check(status.in(200, 201))
    )
    .pause(1)
    .exec(
      http("Login")
        .post("/api/v1/login")
        .body(StringBody("""{"username":"${username}","password":"${password}"}"""))
        .check(
          status.is(200),
          jsonPath("$.token").saveAs("authToken")
        )
    )

  val apiCalls = exec(
    http("Get Single User")
      .get("/api/v1/users/${userId}")
      .header("Authorization", "Bearer ${authToken}")
      .check(status.is(200))
  )
  .pause(1)
  .exec(
    http("Get All Users")
      .get("/api/v1/users")
      .header("Authorization", "Bearer ${authToken}")
      .check(status.in(200, 404))
  )
  .pause(1)
  .exec(
    http("Update User")
      .patch("/api/v1/users/${userId}")
      .header("Authorization", "Bearer ${authToken}")
      .body(StringBody("""{"first_name":"Bench","last_name":"User"}"""))
      .check(status.is(200))
  )

  val scn = scenario("Full API Test")
    .exec(registerAndLogin)
    .pause(1)
    .exec(apiCalls)

  setUp(
    scn.inject(
      rampUsers(1000).during(5.seconds),  // 降低并发用户数，避免数据库压力过大
      constantUsersPerSec(10).during(30.seconds)  // 保持稳定的请求率
    )
  ).protocols(httpProtocol)
    .assertions(
      global.responseTime.max.lt(5000),
      global.successfulRequests.percent.gt(95)
    )
}
