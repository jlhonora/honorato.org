$.getJSON( "http://localhost:7654/api/github", function( data ) {
	var items = [];

	$.each( data, function( key, val ) {

		// Transform each date into readable
		// format with moment.js
		val["readable_date"] = moment(val['created_at']).fromNow()

		// Could've been done with mustache.js, but
		// didn't really add value
		items.push( "<li>" + val['body'] + " " 
			+ "<a href=\"" + val['target']['name_url'] + "\" >" 
			+ val['target']['name'] + "</a><p>" 
			+ val['readable_date'] + "</p></li>" );
	});

	$("#github-table").append(items.join(""));
});
