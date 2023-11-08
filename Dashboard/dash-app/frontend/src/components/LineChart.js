import React, { useRef, useEffect } from 'react';
import * as d3 from 'd3';

const LineChart = ({ data }) => {
  const chartRef = useRef(null);

  useEffect(() => {
    if (data && chartRef.current) {
      const svg = d3.select(chartRef.current);
      svg.selectAll("*").remove(); // Clear previous data

      const width = 900;
      const height = 300;
      const margin = { top: 40, right: 20, bottom: 40, left: 60 }; // Adjusted margin to fit labels

      // Create scales
      const x = d3.scaleTime().domain(d3.extent(data, d => new Date(d.timestamp))).range([margin.left, width - margin.right]);
      const y = d3.scaleLinear().domain([0, d3.max(data, d => d.request_rate)]).nice().range([height - margin.bottom, margin.top]);

      // Create axes
      const xAxis = g => g
        .attr("transform", `translate(0,${height - margin.bottom})`)
        .call(d3.axisBottom(x).ticks(6))
        .attr("font-weight", "bold")
        .attr("font-size", "14px"); // Added font-size as an example of styling

      const yAxis = g => g
        .attr("transform", `translate(${margin.left},0)`)
        .call(d3.axisLeft(y).ticks(height/50))
        .attr("font-weight", "bold")
        .attr("font-size", "14px"); // Added font-size as an example of styling

      // Create line generator
      const line = d3.line()
        .defined(d => !isNaN(d.request_rate))
        .x(d => x(new Date(d.timestamp)))
        .y(d => y(d.request_rate));

      // Line path
      svg.append("path")
        .datum(data)
        .attr("fill", "none")
        .attr("stroke", "#33A1C9")
        .attr("stroke-width", 5)
        .attr("stroke-linejoin", "round")
        .attr("stroke-linecap", "round")
        .attr("d", line);

      // Append axes
      svg.append("g").call(xAxis);
      svg.append("g").call(yAxis);

      // Title
      svg.append("text")
        .attr("x", width / 2)
        .attr("y", margin.top - 15)
        .attr("text-anchor", "middle")
        .attr("font-weight", "bold")
        .attr("font-size", "14px")
        .text("Request Rate");

      // Y Axis label
      svg.append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", 0)
        .attr("x", -(height / 2))
        .attr("dy", "1em")
        .attr("text-anchor", "middle")
        .attr("font-weight", "bold")
        .attr("font-size", "12px")
        .text("Request Rate");

      // X Axis label
      svg.append("text")
        .attr("x", width / 2)
        .attr("y", height + 10)
        .attr("dy", "-1em")
        .attr("text-anchor", "middle")
        .attr("font-weight", "bold")
        .attr("font-size", "12px")
        .text("Time");

      // Hover functionality
      const focus = svg.append("g")
        .attr("class", "focus")
        .style("display", "none");

      focus.append("circle")
        .attr("r", 4.5);

      focus.append("text")
        .attr("x", 9)
        .attr("dy", ".35em");

        svg.append("rect")
        .attr("class", "overlay")
        .attr("width", width)
        .attr("height", height)
        .attr("fill", "none")
        .attr("pointer-events", "all") // This ensures it still listens to mouse events
        .on("mouseover", () => focus.style("display", null))
        .on("mouseout", () => focus.style("display", "none"))
        .on("mousemove", mousemove);
    

        function mousemove(event) {
            const bisect = d3.bisector(d => new Date(d.timestamp)).left;
            const x0 = x.invert(d3.pointer(event, this)[0]);
            const i = bisect(data, x0, 1);
        
            // Check if data points exist
            if (i <= 0 || !data[i - 1] || !data[i]) return;
        
            const d0 = data[i - 1];
            const d1 = data[i];
            const d = x0 - new Date(d0.timestamp) > new Date(d1.timestamp) - x0 ? d1 : d0;
            focus.attr("transform", `translate(${x(new Date(d.timestamp))},${y(d.request_rate)})`);

    // Position the text element relative to the focus group without moving the circle
    focus.select("text")
        .attr("x", 9) // This sets the text offset from the circle
        .attr("y", -10) // This moves the text above the circle
        .text(`${d.request_rate.toFixed(3)} at ${d3.timeFormat("%B %d, %Y %I:%M:%S %p")(new Date(d.timestamp))}`); // Include time in the text
        }
        
    }
  }, [data]);

  return <svg ref={chartRef} width="900" height="300"></svg>;
};

export default LineChart;
