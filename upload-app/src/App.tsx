import { Routes, Route, Link } from "react-router-dom";
import Login from "./pages/Login";
import Upload from "./pages/Upload";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  return (
    <div>
      <nav>
        <Link to="/">Login</Link> |{" "}
        <Link to="/upload">Upload</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Login />} />
        <Route
          path="/upload"
          element={
            <ProtectedRoute>
              <Upload />
            </ProtectedRoute>
          }
        />
      </Routes>
    </div>
  );
}

export default App;
