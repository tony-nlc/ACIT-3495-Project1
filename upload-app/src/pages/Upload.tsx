const Upload = () => {
    return (
        <div>
        <h2>Upload Video</h2>
        <form>
            <div>
            <input type="file" accept="video/*" />
            </div>
            <button type="submit">Upload</button>
        </form>
        </div>
    )
}

export default Upload;