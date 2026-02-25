/**
 * manually created an onChange trigger for the AppsScript in the AppsScript console thingo
 *
 * Alternatively I could sort these rows using the googlesheets go api
 * load all the rows
 * sort them
 * delete all rows
 * set new sorted rows
 * using the googlesheets go api I have
 *
 * column numbers are NOT zero index
 * ie: count starts at 1
 */
function onChange(e) {
    var ss = SpreadsheetApp.getActiveSpreadsheet()
    var sheet = ss.getActiveSheet()

    switch (sheet.getName()) {
        case "Tony's Videos":
            sortTonysVideos(sheet)
            break
        case 'Found Tracks':
            sortFoundTracks(sheet)
            break
        case 'TEST':
            sortTest(sheet)
            break
    }
}

function sortTonysVideos(sheet) {
    var range = sheet.getRange('A2:E')
    range.sort([
        { column: 3, ascending: false }, // published_at
        // { column: 7, ascending: false }, // added_at
    ])
}
function sortFoundTracks(sheet) {
    var range = sheet.getRange('A2:J')
    range.sort([
        { column: 8, ascending: false }, // videopublish date
        { column: 9, ascending: false }, // added_at
    ])
}
function sortTest(sheet) {
    // var range = sheet.getRange('A2:I')
    // range.sort([
    //   { column: 8, ascending: false }, // videopublish date
    //   { column: 9, ascending: false }, // added_at
    // ])
}
