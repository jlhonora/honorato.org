(function () { 
	/*var margin = {top: 0, right: 0, bottom: 0, left: 0},*/
	var margin = {top: 10, right: 10, bottom: 0, left: 35},
		width = 350 - margin.left - margin.right,
		height = 150 - margin.top - margin.bottom;

	var parseDate = d3.time.format("%d-%b-%y").parse;

	var x = d3.time.scale()
		.range([0, width]);

	var y = d3.scale.linear()
		.range([height, 0]);

	var xAxis = d3.svg.axis()
		.scale(x)
		.orient("bottom");

	var yAxis = d3.svg.axis()
		.scale(y)
		.ticks(4)
		.orient("left");

	var line = d3.svg.line()
		.interpolate("bundle") 
		.x(function(d) { return x(d.date); })
		.y(function(d) { return y(d.close); });

	var svg = d3.select("#main-temp-graph").append("svg")
		.attr("width", width + margin.left + margin.right)
		.attr("height", height + margin.top + margin.bottom)
	  .append("g")
		.attr("transform", "translate(" + margin.left + "," + margin.top + ")");

	d3.tsv("static/tempdata.tsv", function(error, data) {
	  data.forEach(function(d) {
		d.date = parseDate(d.date);
		d.close = +d.close;
	  });

	  x.domain(d3.extent(data, function(d) { return d.date; }));
	  y.domain(d3.extent(data, function(d) { return d.close; }));

	/*  svg.append("g")
		  .attr("class", "x axis")
		  .attr("transform", "translate(0," + height + ")")
		  .call(xAxis); */

		  //.attr("transform", "translate(" + width + " ,0)")
	  svg.append("g")
		  .attr("class", "y axis")
		  .call(yAxis)
		.append("text")
		  .attr("dy", "1em")
		  .attr("dx", "4.5em")
		  .style("text-anchor", "end")
		  .text("Temp (C)");

	  svg.append("path")
		  .datum(data)
		  .attr("class", "line-temp")
		  .attr("d", line);

      // Show the element that holds this graph
      $("#sidebar").show();
      $("#iot-wrapper").show();
	
	});
	
})();
