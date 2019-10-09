import React, { useState, useEffect } from 'react'
import * as d3 from "d3";
import LoadingBar from "./LoadingBar"

const BarChartWithLabel = (props) => {
    var errorFound = false
    // error check
    try {
        if (!props.xValueAccessor) {
            throw new Error("Must provide a yValueAccessor function to BarChartWithLabel")
        }
        if(!props.xValueAccessor) {
            throw new Error("Must provice a xValueAccessor function to BarChartWithLabel")
        }
        if(!props.id) {
            throw new Error("Must provide an id to BarChartWithLabel")
        }
    } catch (e) {
        console.error(e, e.stack);
        errorFound = true
    }
    const [chartProperties, setChartProperties] = useState(null)
    const data = props.data
    const defaultSvgHeight = 400;
    const defaultSvgWidth = 800;
    const defaultBarPadding = 5;
    const defaultTopMargin = 30;
    const defaultBottomMargin = 30;
    const defaultLeftMargin = 50;
    const defaultRightMargin = 50;
    const defaultYAxisTicks = 5
    const defaultXAxisPadding = 0.1;
    const defaultBarFillColor = "black"

    const setUpChartProperties = () => {
        // If we have the data but have not yet set chart properties
        if (chartProperties === null && data !== null) {
            var properties = {};
            properties.svgHeight = (props.height) ? props.height : defaultSvgHeight;
            properties.svgWidth = (props.width) ? props.width : defaultSvgWidth;
            properties.margin = {top: (props.topMargin) ? props.topMargin : defaultTopMargin, 
                                right: (props.rightMargin) ? props.rightMargin : defaultRightMargin, 
                                bottom: (props.bottomMargin) ? props.bottomMargin : defaultBottomMargin, 
                                left: (props.leftMargin) ? props.leftMargin : defaultLeftMargin};
            properties.barPadding = (props.barPadding) ? props.barPadding : defaultBarPadding;
            properties.height = properties.svgHeight - properties.margin.top - properties.margin.bottom
            properties.width = properties.svgWidth - properties.margin.left - properties.margin.right
            properties.yAxisNumTicks = (props.yAxisNumTicks) ? props.yAxisNumTicks : defaultYAxisTicks
            properties.xAxisPadding = (props.xAxisPadding) ? props.xAxisPadding : defaultXAxisPadding
            properties.barFillColor = (props.barFillColor) ? props.barFillColor : defaultBarFillColor
            properties.yScale = d3.scaleLinear()
            .rangeRound([properties.height, 0])
            .domain([0, d3.max(data, function(d){
                return props.yValueAccessor(d);
              })]);

            properties.xScale = d3.scaleBand()
            .rangeRound([0, properties.width])
            .domain(data.map(function(d){
                return props.xValueAccessor(d);
            }))
            .padding(properties.xAxisPadding);

            properties.yAxis = d3.axisLeft(properties.yScale)
            .ticks(properties.yAxisNumTicks);
 
            properties.xAxis = d3.axisBottom(properties.xScale)

            setChartProperties(properties)
        }
    }
    const drawChart = () => {
        const chart = d3.select("#"+props.id)
        .append("svg")
        .attr("width", chartProperties.svgWidth)
        .attr("height", chartProperties.svgHeight)
        .append('g')
        .attr('transform', 'translate(' + chartProperties.margin.left + ',' + chartProperties.margin.top + ')');

        chart.append("g")
            .attr("class", "axis axis--x")
            .attr("transform", "translate(0," + chartProperties.height + ")")
            .call(chartProperties.xAxis);

        chart.append("g")
          .attr("class", "axis axis--y")
          .call(chartProperties.yAxis)
          .append("text")
          .attr("transform", "rotate(-90)")
          .attr("y", 6)
          .attr("dy", "0.71em")
          .attr("text-anchor", "end")
          .text("Frequency");

          chart.selectAll(".bar")
          .data(data)
          .enter().append("rect")
            .attr("class", "bar")
            .attr("x", function(d) { return chartProperties.xScale(props.xValueAccessor(d)); })
            .attr("y", function(d) { return chartProperties.yScale(props.yValueAccessor(d)); })
            .attr("width", chartProperties.xScale.bandwidth())
            .attr("height", function(d) { return chartProperties.height - chartProperties.yScale(props.yValueAccessor(d)); })
            .attr("fill", chartProperties.barFillColor);
    };
    useEffect(() => {
        if (data && chartProperties !== null) {
            drawChart();
        } else if (data){
            setUpChartProperties();
        }
    });

    if (errorFound) {
        return null;
    }
    if (!data) {
        return <LoadingBar size={(props.width) ? props.width/4 : defaultSvgWidth/4}/>
    }

    return (
        <div id={props.id}></div>
    ); 
};

export default BarChartWithLabel;