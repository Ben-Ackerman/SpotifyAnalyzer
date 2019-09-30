import React, { useState, useEffect } from 'react'
import BarChartWithLabel from "./BarChartWithLabel"

const UserWordCountChart = (props) => {
    const [data, setData] = useState(null)

    const setupData = async () => {
        if (data === null) {
            setData([{"word":"yeah","count":117},{"word":"love","count":70},{"word":"gonna","count":56},{"word":"high","count":52},{"word":"night","count":51},{"word":"time","count":46},{"word":"'cause","count":45},{"word":"wanna","count":44},{"word":"life","count":43},{"word":"baby","count":42},{"word":"back","count":40},{"word":"play","count":38},{"word":"make","count":38},{"word":"mama","count":36},{"word":"heart","count":34},{"word":"hey","count":33},{"word":"rock","count":32},{"word":"feelin'","count":27},{"word":"find","count":26},{"word":"long","count":24}]);
            // fetch("/v1/my/wordcount", {
            //     method: 'GET',
            //     credentials: 'include'
            // }).then(res => 
            // {
            //     console.log(res.status)
            //     res.json()
            // })
            // .then(json => setData(json["result"])
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