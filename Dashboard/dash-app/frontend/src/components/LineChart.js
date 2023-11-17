import React, { useRef, useEffect } from 'react';
import * as d3 from 'd3';

const LineChart = ({ dataSets, title }) => {
  const chartRef = useRef(null);

  useEffect(() => {
    if (dataSets.length > 0 && chartRef.current) {
      const svg = d3.select(chartRef.current);
      svg.selectAll("*").remove(); // Clear previous data

      const width = 1800;
      const height = 300;
      const margin = { top: 40, right: 150, bottom: 40, left: 40 }; // Adjusted right margin for legend

      // Calculate the overall extent of timestamps across all datasets
      const timeExtents = dataSets.reduce((acc, dataSet) => {
        const extent = d3.extent(dataSet.data, d => new Date(d.timestamp));
        return [d3.min([acc[0], extent[0]]), d3.max([acc[1], extent[1]])];
      }, [Infinity, -Infinity]);

      const x = d3.scaleTime()
        .domain(timeExtents)
        .range([margin.left, width - margin.right]);
      
      const maxY = d3.max(dataSets.map(ds => d3.max(ds.data, d => d.value)));
      const y = d3.scaleLinear()
        .domain([0, maxY])
        .nice()
        .range([height - margin.bottom, margin.top]);

      // Create axes
      const xAxis = g => g
        .attr("transform", `translate(0,${height - margin.bottom})`)
        .call(d3.axisBottom(x).ticks(width / 80).tickSizeOuter(0));

      const yAxis = g => g
        .attr("transform", `translate(${margin.left},0)`)
        .call(d3.axisLeft(y));

      svg.append("g").call(xAxis);
      svg.append("g").call(yAxis);

      // Generate a line for each dataset
      dataSets.forEach((dataSet, i) => {
        const line = d3.line()
          .defined(d => !isNaN(d.value))
          .x(d => x(new Date(d.timestamp)))
          .y(d => y(d.value));

        svg.append("path")
          .datum(dataSet.data)
          .attr("fill", "none")
          .attr("stroke", d3.schemeCategory10[i % 10])
          .attr("stroke-width", 2)
          .attr("d", line);
      });

      // Add title
      svg.append("text")
        .attr("x", width / 2)
        .attr("y", margin.top / 2)
        .attr("text-anchor", "middle")
        .style("font-size", "16px")
        .style("text-decoration", "underline")
        .text(title);

      // Add a legend
      const legendSpacing = 20;
      const legend = svg.append("g")
        .attr("class", "legend")
        .attr("transform", `translate(${width - margin.right + 20},${margin.top})`);

      dataSets.forEach((dataSet, i) => {
        const legendItem = legend.append("g")
          .attr("transform", `translate(0, ${i * legendSpacing})`);

        legendItem.append("rect")
          .attr("width", 10)
          .attr("height", 10)
          .attr("fill", d3.schemeCategory10[i % 10]);

        legendItem.append("text")
          .attr("x", 15)
          .attr("y", 10)
          .attr("text-anchor", "start")
          .style("font-size", "12px")
          .text(dataSet.type);
      });

      // Ensure the legend is within the bounds of the chart
      const legendHeight = dataSets.length * legendSpacing;
      if (height - margin.bottom - legendHeight < 0) {
        legend.attr("transform", `translate(${width - margin.right + 20},${height - margin.bottom - legendHeight})`);
      }
      // // Hover functionality
      // const focus = svg.append("g")
      //   .attr("class", "focus")
      //   .style("display", "none");

      // focus.append("circle")
      //   .attr("r", 4.5);

      // focus.append("text")
      //   .attr("x", 9)
      //   .attr("dy", ".35em");

      // svg.append("rect")
      //   .attr("class", "overlay")
      //   .attr("width", width - margin.left - margin.right)
      //   .attr("height", height - margin.top - margin.bottom)
      //   .attr("transform", `translate(${margin.left},${margin.top})`)
      //   .attr("fill", "none")
      //   .attr("pointer-events", "all")
      //   .on("mouseover", () => focus.style("display", null))
      //   .on("mouseout", () => focus.style("display", "none"))
      //   .on("mousemove", mousemove);

      // function mousemove(event) {
      //   const x0 = x.invert(d3.pointer(event)[0]);
      //   dataSets.forEach(dataSet => {
      //     const bisect = d3.bisector(d => new Date(d.timestamp)).left;
      //     const i = bisect(dataSet.data, x0, 1);
      //     const d0 = dataSet.data[i - 1];
      //     const d1 = dataSet.data[i];
      //     const d = x0 - new Date(d0.timestamp) > new Date(d1.timestamp) - x0 ? d1 : d0;
          
      //     focus.select("circle")
      //       .attr("transform", `translate(${x(new Date(d.timestamp))},${y(d.value)})`);
          
      //     focus.select("text")
      //       .attr("transform", `translate(${x(new Date(d.timestamp))},${y(d.value)})`)
      //       .text(`at ${d3.timeFormat("%B %d, %Y %I:%M:%S %p")(new Date(d.timestamp))}`);
      //   });
      // }
      
    }
  }, [dataSets, title]); // Only re-run the effect if dataSets or title changes

  return <svg ref={chartRef} width="1800" height="300"></svg>;
};

export default LineChart;
