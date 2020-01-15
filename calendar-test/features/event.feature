# file: features/event.feature
Feature: Creation of events, retrieving of notification, retrieving of dayly, weekly and monthly schedules 
    As API client of Calendar Schedule and Notification service. 
    In order to manage of schedule appropriately user may:
        1. create events 
        2. retrieve notifications
        3. retrieve dayly, weekly and monthly schedules
 
    Scenario: User creates today event and receives notification
        When User creates today event:
            | title                | startTime | endTime | notifyBefore | userID |
            | Test of Notification | +30       | +45     | 30           | 1      |
        Then User receives notification with title "Test of Notification" 

    Scenario: User creates event for specified day and check daily schedule
        When User creates day event:
            | title                             | startTime            | endTime              | notifyBefore | userID |
            | Test of Day List 2 Mar 2020 12:00 | 2020-03-02T12:00:00Z | 2020-03-02T12:30:00Z | 0            | 1      |
        Then User's daily schedule contains an event that has been created:
            | title                             | startTime            | endTime              | notifyBefore | userID |
            | Test of Day List 2 Mar 2020 12:00 | 2020-03-02T12:00:00Z | 2020-03-02T12:30:00Z | 0            | 1      |

    Scenario: User creates events for week specified in settings and check his weekly schedule
        When User creates events for week:
            | title                                     | startTime            | endTime              | notifyBefore | userID |
            | Test of weekly schedule 6 Apr 2020 12:00  | 2020-04-06T12:00:00Z | 2020-04-06T12:30:00Z | 0            | 1      |
            | Test of weekly schedule 8 Apr 2020 15:00  | 2020-04-08T15:00:00Z | 2020-04-08T15:30:00Z | 0            | 1      |
        Then User's weekly schedule contains all events that has been created:
            | title                                     | startTime            | endTime              | notifyBefore | userID |
            | Test of weekly schedule 6 Apr 2020 12:00  | 2020-04-06T12:00:00Z | 2020-04-06T12:30:00Z | 0            | 1      |
            | Test of weekly schedule 8 Apr 2020 15:00  | 2020-04-08T15:00:00Z | 2020-04-08T15:30:00Z | 0            | 1      |

    Scenario: User creates events for month specified in settings and check his monthly schedule
        When User creates events for month:
            | title                                      | startTime            | endTime              | notifyBefore | userID |
            | Test of monthly schedule 1 May 2020 9:00   | 2020-05-01T09:00:00Z | 2020-05-01T09:30:00Z | 0            | 1      |
            | Test of monthly schedule 15 May 2020 12:00 | 2020-05-15T12:00:00Z | 2020-05-15T12:30:00Z | 0            | 1      |
            | Test of monthly schedule 30 May 2020 18:00 | 2020-05-30T18:00:00Z | 2020-05-30T18:30:00Z | 0            | 1      |
        Then User's monthly schedule contains all events that has been created:
            | title                                      | startTime            | endTime              | notifyBefore | userID |
            | Test of monthly schedule 1 May 2020 9:00   | 2020-05-01T09:00:00Z | 2020-05-01T09:30:00Z | 0            | 1      |
            | Test of monthly schedule 15 May 2020 12:00 | 2020-05-15T12:00:00Z | 2020-05-15T12:30:00Z | 0            | 1      |
            | Test of monthly schedule 30 May 2020 18:00 | 2020-05-30T18:00:00Z | 2020-05-30T18:30:00Z | 0            | 1      |

