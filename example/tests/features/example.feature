Feature: example
  In order to test my example application
  I need to be able to run it

  Scenario: Graceful shutdown
    Given Quick exit signal delivered
    When application exits with status 0

  Scenario: Echo
    * Call echo "/echo/json?a=foo&b=bar" returns:
    """
    {
        "a": "foo",
        "b": "bar"
    }
    """
    When application exits with status 0
