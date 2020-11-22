Feature: find closest data neighbours
  In order to get closest data neighbours
  As an API user
  I need to be able to add some data to the service
  and then query it with the data that I want to find closest neighbours around

  Scenario: should get closest data by ID
    When I add a set of data to the service:
      | 0.0,0.0,0.0 | REF_0 |
      | 0.0,0.0,1.0 | REF_1 |
      | 0.0,1.0,1.0 | REF_2 |
      | 1.0,1.0,1.0 | REF_3 |
      | 1.0,1.0,0.1 | REF_4 |
      | 1.0,0.0,0.1 | REF_5 |

    When I send a find request with expected size 0 and body:
      | 0.0,0.0,0.0 | REF_0 |
    Then the response code should be 400

    When I send a find request with expected size 1 and body:
      | 0.0,0.0,0.0 | REF_0 |
    Then the response code should be 200
    And the response contains the following elements:
      | 0.0,0.0,0.0 | REF_0 |

    When I send a find request with expected size 3 and body:
      | 0.0,0.0,0.0 | REF_0 |
    Then the response code should be 200
    And the response contains the following elements:
      | 0.0,0.0,0.0 | REF_0 |
      | 0.0,0.0,1.0 | REF_1 |
      | 1.0,0.0,0.1 | REF_5 |

  Scenario: should add data to the repository
    When I add data to the service:
      | 0.0,0.0,0.0 | REF_0 |
    Then the response code should be 201
    Then the response has header "location" matching regex "^/data/(.*)$"
