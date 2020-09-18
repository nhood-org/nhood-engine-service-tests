Feature: find closest data neighbours
  In order to get closest data neighbours
  As an API user
  I need to be able to add some data to the service
  and then query it with the data that I want to find closest neighbours around

  Scenario: should get closest data by ID
    When I add a set of data to the service:
      | ID_0      | 0.0,0.0,0.0   |
      | ID_1      | 0.0,0.0,1.0   |
      | ID_2      | 0.0,1.0,1.0   |
      | ID_3      | 1.0,1.0,1.0   |
      | ID_4      | 1.0,1.0,0.1   |
      | ID_5      | 1.0,0.0,0.1   |

    When I send a find request with expected size 0 and body:
      | ID_0      | 0.0,0.0,0.0   |
    Then the response code should be 400

    When I send a find request with expected size 1 and body:
      | ID_0      | 0.0,0.0,0.0   |
    Then the response code should be 200
    And the response contains the following elements:
      | ID_0      | 0.0,0.0,0.0   |

    When I send a find request with expected size 3 and body:
      | ID_0      | 0.0,0.0,0.0   |
    Then the response code should be 200
    And the response contains the following elements:
      | ID_0      | 0.0,0.0,0.0   |
      | ID_1      | 0.0,0.0,1.0   |
      | ID_5      | 1.0,0.0,0.1   |
