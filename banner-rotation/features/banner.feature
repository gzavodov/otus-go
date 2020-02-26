# file: features/banner.feature
Feature: Creation of banners, banner slots, banner social groups. Binding banners to slots. Retrieving banner for show in specified slot. Registration banner click events for specified slot and social group. 
    As API client of Banner Rotation Service. 
    In order to manage of banners collection appropriately client may:
        1. create banner social groups
        2. create banner slots
        3. create banners and bind them to slots
        4. retrieve banner for show in specified slot
        5. register banner click events for specified slot and social group
    Scenario: Client creates banner social group for further usage
        When Client creates following banner social groups:
            | caption   |
            | teenagers |
            | adults    |
            | elderly   |
        Then Recently created banner social groups are available for using:
            | caption   |
            | teenagers |
            | adults    |
            | elderly   |

    Scenario: Client creates banner slot for further usage
        When Client creates following banner slots:
            | caption    |
            | top-center |
            | top-right  |
            | top-left   |
        Then Recently created banner slots are available for using:
            | caption    |
            | top-center |
            | top-right  |
            | top-left   |

    Scenario: Client creates banners and bind them to slots for subsequent display
        When Client creates following banners and bind to specified slots:
            | banner        | slot       |
            | top-center-1  | top-center |
            | top-center-2  | top-center |
            | top-center-3  | top-center |
            | top-center-4  | top-center |
            | top-center-5  | top-center |
            | top-center-6  | top-center |
            | top-center-7  | top-center |
            | top-center-8  | top-center |
            | top-center-9  | top-center |
            | top-center-10 | top-center |
            | top-right-1   | top-right  |
            | top-right-2   | top-right  |
            | top-right-3   | top-right  |
            | top-right-4   | top-right  |
            | top-right-5   | top-right  |
            | top-right-6   | top-right  |
            | top-right-7   | top-right  |
            | top-right-8   | top-right  |
            | top-right-9   | top-right  |
            | top-right-10  | top-right  |
            | top-left-1    | top-left   |
            | top-left-2    | top-left   |
            | top-left-3    | top-left   |
            | top-left-4    | top-left   |
            | top-left-5    | top-left   |
            | top-left-6    | top-left   |
            | top-left-7    | top-left   |
            | top-left-8    | top-left   |
            | top-left-9    | top-left   |
            | top-left-10   | top-left   |
        Then Recently created banner slots are available for using and bound to appropriate slots:
            | banner        | slot       |
            | top-center-1  | top-center |
            | top-center-2  | top-center |
            | top-center-3  | top-center |
            | top-center-4  | top-center |
            | top-center-5  | top-center |
            | top-center-6  | top-center |
            | top-center-7  | top-center |
            | top-center-8  | top-center |
            | top-center-9  | top-center |
            | top-center-10 | top-center |
            | top-right-1   | top-right  |
            | top-right-2   | top-right  |
            | top-right-3   | top-right  |
            | top-right-4   | top-right  |
            | top-right-5   | top-right  |
            | top-right-6   | top-right  |
            | top-right-7   | top-right  |
            | top-right-8   | top-right  |
            | top-right-9   | top-right  |
            | top-right-10  | top-right  |
            | top-left-1    | top-left   |
            | top-left-2    | top-left   |
            | top-left-3    | top-left   |
            | top-left-4    | top-left   |
            | top-left-5    | top-left   |
            | top-left-6    | top-left   |
            | top-left-7    | top-left   |
            | top-left-8    | top-left   |
            | top-left-9    | top-left   |
            | top-left-10   | top-left   |

    Scenario: Client makes query about banner show
        When Client makes query about banner show for following slots and social groups:
            | slot       | group     |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-center | teenagers |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-right  | adults    |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-left   | elderly   |
            | top-left   | elderly   |
        Then Client receives notification about banner show

    Scenario: Client registers banner click events
        When Client registers banner click event for banner selected on previous step:
        Then Client receives notification about banner click