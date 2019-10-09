import React from 'react'
import UserWordCountChart from './UserWordCountChart'

const UserWordCount = () => {
    return (
        <div>
            <h1>Your Spotify Charts</h1>
            <div>
                <h3>Your top 50 songs word count in the short term</h3>
                <UserWordCountChart length="short"/>
            </div>
            <div>
                <h3>Your top 50 songs word count in the medium term</h3>
                <UserWordCountChart length="medium"/>
            </div>
            <div>
                <h3>Your top 50 songs word count in the long term</h3>
                <UserWordCountChart length="long"/>
            </div>
        </div>
    );
};

export default UserWordCount;