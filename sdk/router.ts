import { Router, Request, Response } from 'express';

interface CandidateBody {
  session_id: string;
  candidate: RTCIceCandidate;
}

const candidateStore = new Map<string, RTCIceCandidate[]>();

/**
 * Creates an Express router for storing and fetching ICE candidates
 * @returns Express.Router
 */
export function createCandidateRouter(): Router {
  const router = Router();

  router.post('/candidate', (req: any, res: any) => {
    try {
      const { session_id, candidate } = req.body;
      if (!session_id || !candidate) {
        return res.status(400).json({ error: 'Missing session_id or candidate' });
      }
      if (!candidateStore.has(session_id)) {
        candidateStore.set(session_id, []);
      }
      candidateStore.get(session_id)!.push(candidate);

      res.status(200).json({ message: 'ICE candidate received successfully' });
    } catch (error) {
      console.error('Error processing ICE candidate:', error);
      res.status(500).json({ error: 'Internal server error' });
    }
  });

  router.get('/candidate/:session_id', (req: Request<{ session_id: string }>, res: Response) => {
    const { session_id } = req.params;
    const sessionCandidates = candidateStore.get(session_id) || [];
    res.status(200).json(sessionCandidates);
  });

  return router;
}
