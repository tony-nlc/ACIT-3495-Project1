import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import Login from "./pages/Login";
import Streaming from "./pages/Streaming";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  return (
    <BrowserRouter>
      <div>
        <nav>
          <Link to="/">Login</Link> |{" "}
          <Link to="/stream">Streaming</Link>
        </nav>

        <Routes>
          <Route path="/" element={<Login />} />
          <Route
            path="/stream"
            element={
              <ProtectedRoute>
                <Streaming />
              </ProtectedRoute>
            }
          />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
