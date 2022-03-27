Feature: Picture resizing
    Scenario: Duyên wants to receive a resized picture
        """ Offtopic:
                Duyên is a vietnamese girl firstname. It means destiry.
                It is interesting that "D" reads as english "Z",
                and "Đ" reads as english "D".
        """
        Given Duyên selects the size '<width>x<height>'
        And a link to existen picture of a bigger size of type 'image/jpeg'
        When Duyên calls an endpoint
        Then she receives an image of a content type 'image/jpeg'
        And the size of the image is '<width>x<height>'
        Examples:
        | <width> | <height> |
        | 256     | 128      |
        | 128     | 128      |
        | 1       | 1        |

    Scenario: Duyên provides an invalid size
        Given Duyên selects the size <size>
        And a link to an existen picture of a bigger size
        When Duyên calls an endpoint
        Then she receives an error
        Examples:
        | <size> |
        | 0x0    |
        | 0x1    |
        | ax10   |
        | 10,10  |
        | -1x10  |
        | xxx    |
        | 1x1x1  |
    
    Scenario: Duyên provides an invalid link
        Given Duyên selects the size 256x256
        And a link to an unexistent picture
        When Duyên calls an endpoint
        Then she receives an error

    Scenario: Duyên provides not jpeg image
        Given Duyên selects the size 256x256
        And a link to an existent plain/text file
        When Duyên calls an endpoint
        Then she receives an error
