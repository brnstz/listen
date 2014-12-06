$(document).ready(function() {
    $.getJSON("api/shows.json", function(data) {
        var shows = $("div#shows");

        var maxBands = 2;
        var maxTracks = 3;

        for (var i = 0; i < data.length; i++) {
            // Create a row for this show
            var newRow = $("<div>", {class: "row"});
          
            // Add the date column
            var colDate = $("<div>", {class: "col-md-2"});
            $(colDate).text(data[i].starts_formatted);
            newRow.append(colDate);

            // Add the venue column
            var colVenue = $("<div>", {class: "col-md-4"});
            $(colVenue).text(data[i].venue.name);
            newRow.append(colVenue);

            // Add bands and tracks
            for (var j = 0; (j < data[i].bands) && (j < maxBands); j++) {
                console.log("hello");

                var colBand = $("<div>", {class: "col-md-2"});
                $(colBand).text(data[i].bands[j].name);

                for (var k = 0; (k < data[i].bands[j].tracks) && (k < maxTracks); k++) {
                    $(colBand).append(data[i].bands[j].tracks[k].html);
                }

                newRow.append(colBand);
            }

            // Add row to main div of shows
            shows.append(newRow);
        }
    });
});
