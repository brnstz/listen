$(document).ready(function() {
    $.getJSON("api/shows.json", function(data) {
        var shows = $("div#shows");

        var maxTracks = 3;

        for (var i = 0; i < data.length; i++) {
            // Create a row for this show
            var newRow = $("<div>", {"class": "row"});
          
            // Add the date column
            var colDate = $("<div>", {"class": "col-md-3"});
            $(colDate).text(data[i].starts_formatted);
            newRow.append(colDate);

            // Add the venue column
            var colVenue = $("<div>", {"class": "col-md-3"});
            $(colVenue).text(data[i].venue.name);
            newRow.append(colVenue);

            // Create a column for the bands
            var colBand = $("<div>", {"class": "col-md-6"});

            // Add bands and tracks
            for (var j = 0; data[i].bands != null && j < data[i].bands.length; j++) {
                var band = data[i].bands[j];

                if (band.tracks != null && band.tracks.length > 0) {

                    // If band has tracks, then make a link
                    // Create an anchor tag for this band 
                    var bandAnchor = $("<a>", {
                        "data-toggle":    "popover", 
                        "data-placement": "bottom",
                    });
                    $(bandAnchor).text(band.name);

                    // Append the band
                    $(colBand).append(bandAnchor);

                    // Create the track list
                    var trackPop = "";
                    for (var k = 0; k < band.tracks.length && k < maxTracks; k++) {
                        trackPop += band.tracks[k].html;
                    }

                    $(bandAnchor).popover({
                        "html":     true, 
                        "content":  trackPop,
                        "trigger":  "hover",
                        "container": bandAnchor
                    });
           
                } else {
                    // If no tracks, then just list name of the band
                    $(colBand).append(band.name);
                }

                // Add a comma if it's not the last band
                if (j+1 < data[i].bands.length) {
                    $(colBand).append(", ");
                }
  
            }

            // Append the band column to our row
            newRow.append(colBand);

            // Add row to main div of shows
            shows.append(newRow);
        }
    });
});
