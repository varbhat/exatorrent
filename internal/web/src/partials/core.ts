import { writable, get } from 'svelte/store';
import type { Writable } from 'svelte/store';
import slocation from 'slocation';

export interface DlObject {
  infohash: string;
  name?: string;
  bytescompleted?: number;
  bytesmissing?: number;
  length?: number;
  state: string;
  seeding?: boolean;
}

interface RespObject {
  type: string;
  state: string;
  infohash?: string;
  message: string;
}

interface FsFile {
  name: string;
  path: string;
  size: number;
  isdir: boolean;
}

interface TorrentFile {
  bytescompleted: number;
  displaypath: string;
  length: number;
  offset: number;
  path: string;
  priority: number;
}

interface DevInfo {
  arch: string;
  numbercpu: number;
  cpumodel: string;
  hostname: string;
  platform: string;
  os: string;
  totalmem: number;
  goversion: string;
  startedat: string;
}
interface DevStats {
  cpucycles: number;
  diskfree: number;
  diskpercent: number;
  mempercent: number;
  gomem: number;
  gomemsys: number;
  goroutines: number;
}

interface TorrentStats {
  TotalPeers: number;
  PendingPeers: number;
  ActivePeers: number;
  ConnectedSeeders: number;
  HalfOpenPeers: number;
  BytesWritten: number;
  BytesWrittenData: number;
  BytesRead: number;
  BytesReadData: number;
  BytesReadUsefulData: number;
  ChunksWritten: number;
  ChunksRead: number;
  ChunksReadUseful: number;
  ChunksReadWasted: number;
  MetadataChunksRead: number;
  PiecesDirtiedGood: number;
  PiecesDirtiedBad: number;
}

interface UserConn {
  username: string;
  isadmin: boolean;
  contime: string;
}

interface Userrepr {
  Username: string;
  Token: string;
  UserType: number;
  CreatedAt: string;
}

interface TorcSettings {
  TUploadRateLimiter: number;
  TDownloadRateLimiter: number;
  DisableLocalCache: boolean;
  OnlineCacheURL: string;
  TrackerRefresh: number;
  TrackerListURLs: string[];
  DisAllowTrackersUser: boolean;
  DisAllowTrackersCache: boolean;
  GlobalSeedRatio: number;
  SRRefresh: number;
  DontRemoveCacheInfo: boolean;
}

interface DiskStats {
  total: number;
  free: number;
  used: number;
  usedPercent: number;
}

export let socket: WebSocket;
export const isSignedIn = writable(false);
export const isDisConnected = writable(false);
export const filepagediscon = writable(false);

export const dontstart = writable('false');
export const isAdmin = writable(false);
export const downloadslist: Writable<{ has: boolean; data: DlObject[] }> = writable({
  has: false,
  data: []
});
export const torrentinfo: Writable<DlObject> = writable({} as DlObject);
export const istrntlocked: Writable<boolean> = writable(false);
export const resplist: Writable<{ has: boolean; data: RespObject[] }> = writable({
  has: false,
  data: []
});

export const torrentstats: Writable<TorrentStats> = writable({} as TorrentStats);

export const torcstatus: Writable<Object> = writable({});
export const machinfo: Writable<DevInfo> = writable({} as DevInfo);
export const machstats: Writable<DevStats> = writable({} as DevStats);

export const fsdirinfo: Writable<FsFile[]> = writable([]);
export const torrentfiles: Writable<TorrentFile[]> = writable([]);

export const fileviewpath: Writable<string> = writable('');
export const fileviewinfohash: Writable<string> = writable('');

export const fsfileinfo: Writable<FsFile> = writable({} as FsFile);
export const torrentfileinfo: Writable<TorrentFile> = writable({} as TorrentFile);
export const adminmode: Writable<boolean> = writable(false);

export const usersfortorrent: Writable<string[]> = writable([] as string[]);
export const torrentsforuser: Writable<string[]> = writable([] as string[]);

export const userconnlist: Writable<UserConn[]> = writable([] as UserConn[]);
export const userlist: Writable<Userrepr[]> = writable([] as Userrepr[]);
export const engconfig: Writable<TorcSettings> = writable({} as TorcSettings);

export const torctime: Writable<{ addedat: string; startedat: string }> = writable({ addedat: '', startedat: '' });

export const diskstats: Writable<DiskStats> = writable({} as DiskStats);

export const nooftrackersintrackerdb: Writable<number> = writable(0);

export const hasMachinfo: Writable<boolean> = writable(false);

export const terrormsg: Writable<{ has: boolean; msg: string }> = writable({ has: true, msg: '' });
export const versionstr: Writable<string> = writable('');
export const versionchecked: Writable<boolean> = writable(false);

let curobj: Location;
let un = '';
let pw = '';
let firsttimecon = true;

let wonopenfn = () => {
  firsttimecon = false;
  isDisConnected.set(false);
  localStorage.getItem('exausertype') === 'admin' ? isAdmin.set(true) : isAdmin.set(false);
  curobj = get(slocation);
  if (curobj?.pathname === '/signin') {
    slocation.goto('/');
  }
};

let wonclosefn = () => {
  isDisConnected.set(true);
};

let werrorfn = () => {
  if (firsttimecon === false && location.pathname !== '/signin') {
    alert('Error Connecting');
    return;
  }
  try {
    fetch('/api/auth', {
      method: 'POST',
      body: JSON.stringify({ data1: un, data2: pw })
    })
      .then((res) => {
        if (res.status >= 200 && res.status <= 299) {
          return res.json();
        } else {
          alert('Error Authenticating');
          throw new Error('Error Authenticating');
        }
      })
      .then((res) => {
        localStorage.setItem('exasession', res?.session);
        localStorage.setItem('exausertype', res?.usertype);

        if (!(localStorage.getItem('dontstart') == undefined)) {
          localStorage.setItem('dontstart', 'false');
        } else {
          dontstart.set(localStorage.getItem('dontstart'));
        }

        socket = new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/api/socket');
        socket.onopen = wonopenfn;
        socket.onmessage = SocketHandler;
        socket.onclose = wonclosefn;
        socket.onerror = werrorfn;
        isSignedIn.set(true);
      });
  } catch (err) {
    console.log(err);
    SignOut();
    return;
  }
};

export let Connect = () => {
  un = localStorage.getItem('exausername');
  pw = localStorage.getItem('exapassword');

  if (un != '' && un != undefined && un != null) {
    if (pw != '' && pw != undefined && pw != null) {
      console.log('Signing In');
    } else {
      slocation.goto('/signin');
      return;
    }
  } else {
    slocation.goto('/signin');
    return;
  }

  if (!(un.length > 5) || !(pw.length > 5)) {
    alert('Invalid Credentials');
    return;
  }

  if (socket != null || socket != undefined) {
    socket?.close();
  }

  socket = new WebSocket((window.location.protocol === 'https:' ? 'wss://' : 'ws://') + window.location.host + '/api/socket');

  socket.onopen = wonopenfn;
  socket.onmessage = SocketHandler;
  socket.onclose = wonclosefn;
  socket.onerror = werrorfn;
};

export let SocketHandler = (event: MessageEvent) => {
  console.log('On Message:', event.data);
  let msg = JSON.parse(event.data);
  console.log(msg);

  switch (msg.type) {
    case 'resp':
      if (!(msg == null)) {
        if (msg?.state === 'error' || msg?.state === 'success') {
          alert(msg?.message);
        }
        let rl = get(resplist);
        resplist.set({ has: true, data: [msg, ...rl?.data] as RespObject[] });
      } else {
        resplist.set({ has: false, data: [] as RespObject[] });
      }
      break;
    case 'nfn':
      if (!(msg == null)) {
        let rl = get(resplist);
        resplist.set({ has: true, data: [msg, ...rl?.data] as RespObject[] });
      } else {
        resplist.set({ has: false, data: [] as RespObject[] });
      }
      break;
    case 'torrentstream':
      if (!(msg.data == null)) {
        terrormsg.set({ has: false, msg: '' });
        downloadslist.set({ has: true, data: msg.data as DlObject[] });
      } else {
        terrormsg.set({ has: true, msg: 'No Torrents' });
        downloadslist.set({ has: true, data: [] as DlObject[] });
      }
      break;
    case 'torrentinfo':
      if (!(msg.data == null)) {
        terrormsg.set({ has: false, msg: '' });
        torrentinfo.set({
          infohash: msg.data?.infohash,
          name: msg.data?.name,
          bytescompleted: msg.data?.bytescompleted,
          bytesmissing: msg.data?.bytesmissing,
          length: msg.data?.length,
          state: msg.data?.state,
          seeding: msg.data?.seeding
        } as DlObject);
      } else {
        terrormsg.set({ has: true, msg: 'No Torrent Info' });
        torrentinfo.set({} as DlObject);
      }
      break;
    case 'torrentinfostat':
      if (!(msg.data == null)) {
        torctime.set({ addedat: new Date(msg.data?.AddedAt)?.toLocaleString(), startedat: new Date(msg.data?.StartedAt)?.toLocaleString() });
      } else {
        torctime.set({ addedat: '', startedat: '' });
      }
      break;
    case 'torrentstats':
      console.log(msg);
      if (!(msg.data == null)) {
        torrentstats.set(msg.data as TorrentStats);
      } else {
        torrentstats.set({} as TorrentStats);
      }
      break;
    case 'fsdirinfo':
      console.log(msg);
      if (!(msg.data == null)) {
        fsdirinfo.set(msg.data as FsFile[]);
      } else {
        fsdirinfo.set([]);
      }
      break;
    case 'torrentfiles':
      console.log(msg);
      if (!(msg.data == null)) {
        torrentfiles.set(msg.data as TorrentFile[]);
      } else {
        torrentfiles.set([]);
      }
      break;
    case 'torrentmetainfo':
      const linkSource = `data:application/x-bittorrent;base64,${msg?.data}`;
      const downloadLink = document.createElement('a');
      downloadLink.href = linkSource;
      downloadLink.download = `${msg?.infohash}.torrent`;
      downloadLink.click();
      break;
    case 'torrentfileinfo':
      console.log(msg);
      if (!(msg.data == null)) {
        torrentfileinfo.set(msg.data as TorrentFile);
      } else {
        torrentfileinfo.set({} as TorrentFile);
      }
      break;
    case 'fsfileinfo':
      console.log(msg);
      if (!(msg.data == null)) {
        fsfileinfo.set(msg.data as FsFile);
      } else {
        fsfileinfo.set({} as FsFile);
      }
      if (get(filepagediscon) === true) {
        socket?.readyState === WebSocket.OPEN ? socket?.close() : console.log('socket already closed');
      }
      break;
    case 'usersfortorrent':
      console.log(msg);
      if (!(msg.data == null)) {
        usersfortorrent.set(msg.data as string[]);
      } else {
        usersfortorrent.set([]);
      }
      break;
    case 'torrentsforuser':
      console.log(msg);
      if (!(msg.data == null)) {
        torrentsforuser.set(msg.data as string[]);
      } else {
        torrentsforuser.set([]);
      }
      break;
    case 'torcstatus':
      console.log(msg);
      if (!(msg.data == null)) {
        torcstatus.set(msg.data as TorrentStats);
      } else {
        torcstatus.set({});
      }
      break;
    case 'torrentlockstate':
      console.log(msg);
      if (!(msg.data == null)) {
        istrntlocked.set(msg.data === true);
      } else {
        istrntlocked.set(false);
      }
      break;
    case 'userconn':
      console.log(msg);
      if (!(msg.data == null)) {
        userconnlist.set(msg.data as UserConn[]);
      } else {
        userconnlist.set([] as UserConn[]);
      }
      break;
    case 'users':
      console.log(msg);
      if (!(msg.data == null)) {
        userlist.set(msg.data as Userrepr[]);
      } else {
        userlist.set([] as Userrepr[]);
      }
      break;
    case 'engconf':
      console.log(msg);
      if (!(msg.data == null)) {
        engconfig.set(msg.data as TorcSettings);
      } else {
        engconfig.set({} as TorcSettings);
      }
      break;
    case 'machinfo':
      console.log(msg);
      if (!(msg.data == null)) {
        hasMachinfo.set(true);
        machinfo.set(msg.data as DevInfo);
      } else {
        machinfo.set({} as DevInfo);
      }
      break;
    case 'machstats':
      console.log(msg);
      if (!(msg.data == null)) {
        machstats.set(msg.data as DevStats);
      } else {
        machstats.set({} as DevStats);
      }
      break;
    case 'diskusage':
      console.log(msg);
      if (!(msg.data == null)) {
        diskstats.set(msg.data as DiskStats);
      } else {
        diskstats.set({} as DiskStats);
      }
      break;
    case 'version':
      console.log(msg);
      if (!(msg.data == null)) {
        versionchecked.set(true);
        versionstr.set(msg.data as string);
      } else {
        versionstr.set('');
      }
      break;
    case 'nooftrackersintrackerdb':
      console.log(msg);
      if (!(msg.data == null)) {
        nooftrackersintrackerdb.set(msg.data as number);
      } else {
        nooftrackersintrackerdb.set(0);
      }
      break;
  }
};

export let SignOut = () => {
  localStorage.removeItem('exausername');
  localStorage.removeItem('exapassword');
  localStorage.removeItem('exasession');
  localStorage.removeItem('exausertype');
  localStorage.removeItem('dontstart');
  // Remove Cookies
  document.cookie.split(';').forEach((c) => {
    document.cookie = c.replace(/^ +/, '').replace(/=.*/, '=;expires=' + new Date().toUTCString() + ';path=/;SameSite=Lax;');
  });
  isDisConnected.set(false);
  slocation.goto('/signin');
};

export let Send = (value: any) => {
  console.log('sending ', value);
  if (socket?.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(value));
  }
};

export let fileSize = (b: number) => {
  let u = 0,
    s = 1024;
  while (b >= s || -b >= s) {
    b /= s;
    u++;
  }
  return (u ? b.toFixed(3) + ' ' : b) + ' KMGTPEZY'[u] + 'B';
};

export let fileType = (filepath: string): string => {
  let ext = filepath.split('.').pop();
  let vidext = ['webm', 'mkv', 'flv', 'vob', 'ogv', 'ogg', 'rrc', 'gifv', 'mng', 'mov', 'avi', 'qt', 'wmv', 'yuv', 'rm', 'asf', 'amv', 'mp4', 'm4p', 'm4v', 'mpg', 'mp2', 'mpeg', 'mpe', 'mpv', 'm4v', 'svi', '3gp', '3g2', 'mxf', 'roq', 'nsv', 'flv', 'f4v', 'f4p', 'f4a', 'f4b'];
  let audext = ['aac', 'aiff', 'ape', 'au', 'flac', 'gsm', 'it', 'm3u', 'm4a', 'mid', 'mod', 'mp3', 'mpa', 'pls', 'ra', 's3m', 'sid', 'wav', 'wma', 'xm'];
  if (vidext.includes(ext)) {
    return 'video';
  } else if (audext.includes(ext)) {
    return 'audio';
  }
  return 'unknown';
};
