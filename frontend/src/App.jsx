import {useState} from 'react';
import './App.css';
import {StartTimer, CancelTimer} from '../wailsjs/go/main/App';

function App() {
    const [action, setAction] = useState('shutdown');
    const [timeValue, setTimeValue] = useState(0);
    const [timeUnit, setTimeUnit] = useState('minutes');
    const [message, setMessage] = useState('');

    const handleGo = () => {
        const time = parseInt(timeValue, 10);
        if (isNaN(time) || time <= 0) {
            setMessage('Please enter a valid time.');
            return;
        }
        StartTimer(action, time, timeUnit)
            .then(() => {
                setMessage(`PC will ${action} in ${time} ${timeUnit}.`);
            })
            .catch(err => {
                setMessage(`Error: ${err}`);
            });
    };

    const handleCancel = () => {
        CancelTimer()
            .then(() => {
                setMessage('Timer cancelled.');
            })
            .catch(err => {
                setMessage(`Error: ${err}`);
            });
    };

    return (
        <>
            <div id="App">
                <div className="container">
                    <h1>PC Power Timer</h1>
                    <div className="controls">
                        <span>PC will</span>
                        <select value={action} onChange={(e) => setAction(e.target.value)}>
                            <option value="shutdown">shutdown</option>
                            <option value="restart">restart</option>
                            <option value="sleep">sleep</option>
                        </select>
                        <span>after</span>
                        <input type="number" value={timeValue} onChange={(e) => setTimeValue(e.target.value)} min="1" />
                        <select value={timeUnit} onChange={(e) => setTimeUnit(e.target.value)}>
                            <option value="seconds">seconds</option>
                            <option value="minutes">minutes</option>
                            <option value="hours">hours</option>
                        </select>
                    </div>
                    <div className="buttons">
                        <button className="btn go" onClick={handleGo}>Go</button>
                        <button className="btn cancel" onClick={handleCancel}>Cancel</button>
                    </div>
                    {message && <div className="message">{message}</div>}
                </div>
            </div>
            <footer className="footer">By <a href="https://github.com/YosriMlik/" target="_blank" rel="noopener noreferrer">Yosri Mlik</a> ©</footer>
        </>
    );
}

export default App;