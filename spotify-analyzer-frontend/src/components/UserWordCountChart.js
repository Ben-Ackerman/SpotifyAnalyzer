import React, { useState, useEffect } from 'react'
import BarChartWithLabel from "./BarChartWithLabel"

const UserWordCountChart = (props) => {
    var errorFound = false
    try {
        // do some crazy stuff
        if (!props.length) {
            throw new Error("Must provide a length to the react function UserWordCountChart")
        }
    } catch (e) {
        console.error(e, e.stack);
        return (<div></div>);
    }

    const [data, setData] = useState(null)
    const setupData = async () => {
        if (data === null) {
            url = "/v1/my/wordcount".concat(props.length)
            fetch(url, {
                method: 'GET',
                credentials: 'include'
            }).then(res => {
                console.log(res.status)
                res.json()
            })
            .then(json => setData(json["result"]));
        }
    }

    const getXValue = (d) => {
        return d.word
    }

    const getYValue = (d) => {
        return d.count
    }

    useEffect(() => {
        if (!data) {
            setupData()
        }
    });

    return (
        <BarChartWithLabel 
            width={props.width} 
            height={props.height} 
            id="usersTopWords" 
            barFillColor="green"
            data={data}
            xValueAccessor={getXValue}
            yValueAccessor={getYValue}
        />
    );
};

export default UserWordCountChart;