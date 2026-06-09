import { useState, useEffect } from 'react';
import { Search, Cpu, HardDrive, Monitor, CheckCircle, XCircle, SearchIcon } from 'lucide-react';
import './App.css';
import { GetHardwareInfo, SearchGame, GetGameRequirements, CheckRequirements } from '../wailsjs/go/main/App';

function App() {
    const [specs, setSpecs] = useState<any>(null);
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<any[]>([]);
    const [selectedGameId, setSelectedGameId] = useState<string | null>(null);
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
            {/* Header / PC Data Section */}
            <header className="header-section glass-panel">
                <div className="logo-area">
                    <h1>Can You Run It? <span>Go</span></h1>
                </div>
                {specs && (
                    <div className="pc-specs-grid">
                        <div className="spec-item">
                            <Cpu size={20} />
                            <div>
                                <span className="label">CPU</span>
                                <strong>{specs.cpu_name}</strong>
                            </div>
                        </div>
                        <div className="spec-item">
                            <Monitor size={20} />
                            <div>
                                <span className="label">GPU</span>
                                <strong>{specs.gpu_name}</strong>
                            </div>
                        </div>
                        <div className="spec-item">
                            <HardDrive size={20} />
                            <div>
                                <span className="label">RAM / Storage</span>
                                <strong>{Math.round(specs.ram_total_mb / 1024)}GB / {specs.disk_total_gb}GB</strong>
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
                            <h2>Requirements</h2>
                            
                            <div className="requirements-grid">
                                <div className="req-box min-req">
                                    <h3>Minimum</h3>
                                    {/* The Steam API often returns HTML string for requirements, so we use dangerouslySetInnerHTML to parse it properly */}
                                    <div dangerouslySetInnerHTML={{__html: gameDetails.min}}></div>
                                </div>
                                <div className="req-box rec-req">
                                    <h3>Recommended</h3>
                                    <div dangerouslySetInnerHTML={{__html: gameDetails.rec}}></div>
                                </div>
                            </div>

                            <div className={`match-section ${gameDetails.match ? 'match-success' : 'match-fail'}`}>
                                {gameDetails.match ? (
                                    <>
                                        <CheckCircle size={32} />
                                        <div className="match-text">
                                            <h3>Your PC Meets the Requirements</h3>
                                            <p>You should be able to run this game smoothly based on minimum specs.</p>
                                        </div>
                                    </>
                                ) : (
                                    <>
                                        <XCircle size={32} />
                                        <div className="match-text">
                                            <h3>Your PC Might Struggle</h3>
                                            <p>Your hardware is below the minimum required settings.</p>
                                        </div>
                                    </>
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
