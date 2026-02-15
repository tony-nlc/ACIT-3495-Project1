import './App.css'
import { Routes, Route } from 'react-router-dom';
import Login from './pages/Login';
import Upload from './pages/Upload';
import { Link } from 'react-router-dom';

function App() {

  return (
      <div>
          <nav>
            <Link to="/">Login</Link> |{" "}
            <Link to="/upload">Upload</Link>
          </nav>

          <Routes>
            <Route path="/" element={<Login />}></Route>
            {/* Add more routes here */}
            <Route path="/upload" element={<Upload />}></Route> 
        </Routes>
      </div>
  )
}

export default App
