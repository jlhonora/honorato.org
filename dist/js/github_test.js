var data;
data = [
	{
		"body": "Pushed to",
		"target": {
		  "name": "jlhonora/airbrake-android",
		  "name_url": "https://github.com/jlhonora/airbrake-android"
		},
		"created_at": "2013-10-09T20:21:35Z"
	},
	{
		"body": "Forked",
		"target": {
		  "name": "loopj/airbrake-android",
		  "name_url": "https://github.com/jlhonora/airbrake-android"
		},
		"created_at": "2013-10-09T20:14:12Z"
	},
	{
		"body": "Pushed to",
		"target": {
		  "name": "jlhonora/lsusb",
		  "name_url": "https://github.com/jlhonora/lsusb"
		},
		"created_at": "2013-09-17T01:03:02Z"
	}];

var items = [];

$.each( data, function( key, val ) {

	// Transform each date into readable
	// format with moment.js
	val["readable_date"] = moment(val['created_at']).fromNow()

	// Could've been done with mustache.js, but
	// didn't really add value
    items.push( "<li>" + val['body'] + " " + "<a href=\"" + val['target']['name_url'] + "\" >" + val['target']['name'] + "</a><p>" + val['readable_date'] + "</p></li>" );
});

$("#github-table").append(items.join(""));
