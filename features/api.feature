Feature: get closest data by ID
  In order to get closest data
  As an API user
  I need to be able to request it by ID

  Scenario: should get closes data by ID
    When I send "GET" request to "/" with ID 1
    Then the response code should be 200
    And the response should match json: 
      """
      {"id": 1}
      """
