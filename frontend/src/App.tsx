import { useState, useEffect } from 'react';
import { Search, Cpu, HardDrive, Monitor, CheckCircle, XCircle, SearchIcon, AlertCircle, Info, Check, X } from 'lucide-react';
import './App.css';
import { GetHardwareInfo, SearchGame, GetGameRequirements, CheckRequirements } from '../wailsjs/go/main/App';

// Strip the leading heading Steam includes in requirement HTML
// e.g. <strong>Minimum:</strong>, <strong>REQUIRED</strong>, etc.
function cleanReqHtml(html: string): string {
    if (!html || html === 'Not specified') return html;
    return html
        // Remove first <strong> that is just a header label (Minimum/Recommended/Required etc.)
        .replace(/^\s*<strong>[^<]*(minimum|recommended|required|requires)[^<]*<\/strong>(\s*<br\s*\/?>\s*)*/gi, '')
        // Remove any bare label line like "Minimum:<br>" at the start
        .replace(/^\s*(minimum|recommended|required|requires)[^<]*(<br\s*\/?>\s*)+/gi, '')
        .trim();
}

function App() {
    const [specs, setSpecs] = useState<any>(null);
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<any[]>([]);
    const [selectedGameId, setSelectedGameId] = useState<string | null>(null);
    const [selectedGameName, setSelectedGameName] = useState<string | null>(null);
    const [history, setHistory] = useState<any[]>([]);
    const [gameDetails, setGameDetails] = useState<any>(null);

    useEffect(() => {
        // Fetch real PC data from backend
        GetHardwareInfo().then(data => {
            setSpecs(data);
        }).catch(err => {
            console.error("Error fetching hardware info:", err);
        });
    }, []);

    const handleSearch = async () => {
        if (!searchQuery.trim()) return;
        try {
            const results = await SearchGame(searchQuery);
            setSearchResults(results || []);
        } catch (err) {
            console.error("Error searching game:", err);
            setSearchResults([]);
        }
    };

    const handleSelectGame = async (game: any) => {
        setSelectedGameId(game.appid); // Note: steam API uses 'appid'
        setSelectedGameName(game.name);
        
        try {
            const reqData = await GetGameRequirements(parseInt(game.appid, 10));
            let matchData = null;
            
            if (specs && reqData) {
                matchData = await CheckRequirements(specs, reqData);
            }

            setGameDetails({
                min: reqData?.pc_requirements?.minimum || "Not specified",
                rec: reqData?.pc_requirements?.recommended || "Not specified",
                match: matchData?.MeetsMinimum || false,
                matchDetails: matchData
            });

            // Update history (last 3)
            setHistory(prev => {
                const newHistory = [game, ...prev.filter(g => g.appid !== game.appid)];
                return newHistory.slice(0, 3);
            });
        } catch (err) {
            console.error("Error fetching game details:", err);
        }
    };

    return (
        <div className="app-container">
            {/* Header: two separate glass panels */}
            <header className="header-section">
                <div className="logo-box glass-panel">
                    <div className="logo-area">
                        <h1>Cyris <span className="logo-sub">— Steam</span></h1>
                    </div>
                </div>

                {specs && (
                    <div className="specs-box glass-panel">
                        <div className="pc-specs-grid">
                            <div className="spec-item">
                                <Cpu size={18} />
                                <div>
                                    <span className="label">CPU</span>
                                    <strong>{specs.cpu_name}</strong>
                                </div>
                            </div>
                            <div className="spec-item">
                                <Monitor size={18} />
                                <div>
                                    <span className="label">GPU</span>
                                    <strong>{specs.gpu_name}</strong>
                                </div>
                            </div>
                            <div className="spec-item">
                                <HardDrive size={18} />
                                <div>
                                    <span className="label">RAM</span>
                                    <strong>{Math.ceil(specs.ram_total_mb / 1024)}GB</strong>
                                </div>
                            </div>
                            <div className="spec-item">
                                <HardDrive size={18} />
                                <div>
                                    <span className="label">Free Storage</span>
                                    <strong>{specs.disk_free_gb}GB</strong>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </header>

            <main className="main-content">
                <div className="left-column">
                    {/* Search Section */}
                    <div className="search-section glass-panel">
                        <h2>Find a Game</h2>
                        <div className="search-bar">
                            <input 
                                type="text" 
                                placeholder="Search games (e.g., Cyberpunk 2077)..." 
                                value={searchQuery}
                                onChange={e => setSearchQuery(e.target.value)}
                                onKeyDown={e => e.key === 'Enter' && handleSearch()}
                                autoComplete="new-password"
                                spellCheck={false}
                            />
                            <button onClick={handleSearch} className="search-btn">
                                <SearchIcon size={20} />
                            </button>
                        </div>

                        {searchResults.length > 0 && (
                            <div className="search-results">
                                {searchResults.map(game => (
                                    <div 
                                        key={game.appid} 
                                        className={`result-item ${selectedGameId === game.appid ? 'active' : ''}`}
                                        onClick={() => handleSelectGame(game)}
                                    >
                                        {game.name}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* History Section */}
                    {history.length > 0 && (
                        <div className="history-section glass-panel mt-4">
                            <h2>Recent Searches</h2>
                            <div className="history-list">
                                {history.map(game => (
                                    <div key={game.appid} className="history-item" onClick={() => handleSelectGame(game)}>
                                        {game.name}
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>

                <div className="right-column">
                    {/* Game Details Section */}
                    {selectedGameId && gameDetails ? (
                        <div className="game-details glass-panel">
                            <h2>Requirements — {selectedGameName}</h2>
                            
                            <div className="requirements-grid">
                                <div className="req-box min-req">
                                    <h3>Minimum</h3>
                                    {/* The Steam API returns HTML strings; cleanReqHtml strips the redundant heading line */}
                                    <div dangerouslySetInnerHTML={{__html: cleanReqHtml(gameDetails.min)}}></div>
                                </div>
                                <div className="req-box rec-req">
                                    <h3>Recommended</h3>
                                    <div dangerouslySetInnerHTML={{__html: cleanReqHtml(gameDetails.rec)}}></div>
                                </div>
                            </div>

                            <div className={`match-section ${(!gameDetails.matchDetails?.CPU?.Found || !gameDetails.matchDetails?.GPU?.Found) ? 'match-unknown' : gameDetails.match ? 'match-good' : 'match-bad'}`}>
                                <div className="match-header">
                                    <div className="match-status-icon">
                                        {(!gameDetails.matchDetails?.CPU?.Found || !gameDetails.matchDetails?.GPU?.Found) ? (
                                            <Info size={32} />
                                        ) : gameDetails.match ? (
                                            <Check size={32} />
                                        ) : (
                                            <X size={32} />
                                        )}
                                    </div>
                                    <div className="match-text">
                                        <h3>System Analysis</h3>
                                        <p>
                                            {!gameDetails.matchDetails?.CPU?.Found || !gameDetails.matchDetails?.GPU?.Found
                                                ? "We couldn't identify your exact CPU/GPU, so we cannot assure you if it will work or not."
                                                : gameDetails.match 
                                                    ? "Your PC Meets the Requirements. You should be able to run this game smoothly."
                                                    : "Your hardware is below the minimum required settings."}
                                        </p>
                                    </div>
                                </div>

                                {gameDetails.matchDetails && (
                                    <div className="missing-details">
                                        <ul className="missing-list">
                                            <li className="icon-list-item">
                                                {gameDetails.matchDetails.CPU?.Found ? 
                                                    (gameDetails.matchDetails.CPU?.Meets ? <Check size={16} color="#facc15" /> : <X size={16} color="#f472b6" />) : 
                                                    <Info size={16} color="#94a1d7" />
                                                }
                                                <span>CPU: {gameDetails.matchDetails.CPU?.Found ? (gameDetails.matchDetails.CPU?.Meets ? "Meets Minimum" : "Below Minimum") : "Not Found"}</span>
                                            </li>
                                            
                                            <li className="icon-list-item">
                                                {gameDetails.matchDetails.GPU?.Found ? 
                                                    (gameDetails.matchDetails.GPU?.Meets ? <Check size={16} color="#facc15" /> : <X size={16} color="#f472b6" />) : 
                                                    <Info size={16} color="#94a1d7" />
                                                }
                                                <span>GPU: {gameDetails.matchDetails.GPU?.Found ? (gameDetails.matchDetails.GPU?.Meets ? "Meets Minimum" : "Below Minimum") : "Not Found"}</span>
                                            </li>

                                            <li className="icon-list-item">
                                                {gameDetails.matchDetails.RAMTotal ? <Check size={16} color="#facc15" /> : <X size={16} color="#f472b6" />}
                                                <span>RAM: {gameDetails.matchDetails.RAMTotal ? "Meets Minimum" : "Below Minimum"}</span>
                                            </li>

                                            <li className="icon-list-item">
                                                {gameDetails.matchDetails.DiskFree ? <Check size={16} color="#facc15" /> : <X size={16} color="#f472b6" />}
                                                <span>Storage: {gameDetails.matchDetails.DiskFree ? "Meets Minimum" : "Below Minimum"}</span>
                                            </li>
                                        </ul>
                                    </div>
                                )}
                            </div>
                        </div>
                    ) : (
                        <div className="empty-state glass-panel">
                            <Search size={48} opacity={0.5} />
                            <h3>No Game Selected</h3>
                            <p>Search and select a game to see if you can run it.</p>
                        </div>
                    )}
                </div>
            </main>
        </div>
    );
}

export default App;
