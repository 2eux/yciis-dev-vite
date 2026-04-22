import { useState } from 'react'
import { 
  Search, Plus, Filter, Download, MoreHorizontal, 
  Clock, Mail, Phone, Users, ChevronLeft, ChevronRight
} from 'lucide-react'

const mockStudents = [
  { id: '1', student_id: 'STU2024001', first_name: 'Aryan', last_name: 'Sharma', gender: 'male', section: 'Class 9-A', roll_number: '001', status: 'active', mobile: '+91 98765 43210', email: 'aryan@dps.edu' },
  { id: '2', student_id: 'STU2024002', first_name: 'Priya', last_name: 'Singh', gender: 'female', section: 'Class 9-A', roll_number: '002', status: 'active', mobile: '+91 98765 43211', email: 'priya@dps.edu' },
  { id: '3', student_id: 'STU2024003', first_name: 'Rahul', last_name: 'Verma', gender: 'male', section: 'Class 9-B', roll_number: '001', status: 'active', mobile: '+91 98765 43212', email: 'rahul@dps.edu' },
  { id: '4', student_id: 'STU2024004', first_name: 'Sneha', last_name: 'Williams', gender: 'female', section: 'Class 9-B', roll_number: '002', status: 'active', mobile: '+91 98765 43213', email: 'sneha@dps.edu' },
  { id: '5', student_id: 'STU2024005', first_name: 'Karan', last_name: 'Patel', gender: 'male', section: 'Class 10-A', roll_number: '001', status: 'inactive', mobile: '+91 98765 43214', email: 'karan@dps.edu' },
  { id: '6', student_id: 'STU2024006', first_name: 'Anjali', last_name: 'Gupta', gender: 'female', section: 'Class 10-A', roll_number: '002', status: 'active', mobile: '+91 98765 43215', email: 'anjali@dps.edu' },
]

const classes = ['All Classes', 'Class 9-A', 'Class 9-B', 'Class 10-A', 'Class 10-B']
const genders = ['All Genders', 'Male', 'Female']

const statusColors: Record<string, string> = {
  active: 'badge-primary',
  inactive: 'badge-secondary',
  transferred: 'badge-outline',
  dropped: 'badge-destructive',
}

export default function Students() {
  const [search, setSearch] = useState('')
  const [selectedClass, setSelectedClass] = useState('All Classes')
  const [selectedGender, setSelectedGender] = useState('All Genders')
  const [selectedStudent, setSelectedStudent] = useState<any>(null)
  const [page, setPage] = useState(1)

  const filteredStudents = mockStudents.filter(student => {
    const matchesSearch = search === '' || 
      student.first_name.toLowerCase().includes(search.toLowerCase()) ||
      student.last_name.toLowerCase().includes(search.toLowerCase()) ||
      student.student_id.toLowerCase().includes(search.toLowerCase())
    const matchesClass = selectedClass === 'All Classes' || student.section === selectedClass
    const matchesGender = selectedGender === 'All Genders' || 
      (selectedGender === 'Male' && student.gender === 'male') ||
      (selectedGender === 'Female' && student.gender === 'female')
    return matchesSearch && matchesClass && matchesGender
  })

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="page-header">
        <div>
          <h1 className="page-title">Students</h1>
          <p className="page-description">Manage student records and information</p>
        </div>
        <div className="flex items-center gap-3">
          <button className="btn btn-secondary">
            <Download className="mr-2 h-4 w-4" />
            Export
          </button>
          <button className="btn btn-primary">
            <Plus className="mr-2 h-4 w-4" />
            Add Student
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Student List */}
        <div className="lg:col-span-3 space-y-4">
          {/* Filters */}
          <div className="card p-4">
            <div className="flex flex-col md:flex-row gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <input
                  type="text"
                  placeholder="Search by name, student ID, or email..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="input pl-11"
                />
              </div>
              <select 
                value={selectedClass}
                onChange={(e) => setSelectedClass(e.target.value)}
                className="flex h-10 w-full md:w-40 items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                {classes.map(c => <option key={c}>{c}</option>}
              </select>
              <select 
                value={selectedGender}
                onChange={(e) => setSelectedGender(e.target.value)}
                className="flex h-10 w-full md:w-36 items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                {genders.map(g => <option key={g}>{g}</option>}
              </select>
              <button className="btn btn-secondary">
                <Filter className="mr-2 h-4 w-4" />
                More Filters
              </button>
            </div>
          </div>

          {/* Student Table */}
          <div className="card p-0 overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-muted">
                  <tr>
                    <th className="table-head">Student</th>
                    <th className="table-head">Class</th>
                    <th className="table-head">Roll No.</th>
                    <th className="table-head">Contact</th>
                    <th className="table-head">Status</th>
                    <th className="table-head"></th>
                  </tr>
                </thead>
                <tbody>
                  {filteredStudents.map((student) => (
                    <tr 
                      key={student.id}
                      className="table-row cursor-pointer"
                      onClick={() => setSelectedStudent(student)}
                    >
                      <td className="table-cell">
                        <div className="flex items-center gap-3">
                          <div className="avatar bg-primary">
                            <span className="text-primary-foreground font-semibold">
                              {student.first_name[0]}{student.last_name[0]}
                            </span>
                          </div>
                          <div>
                            <p className="font-medium">
                              {student.first_name} {student.last_name}
                            </p>
                            <p className="text-xs text-muted-foreground">{student.student_id}</p>
                          </div>
                        </div>
                      </td>
                      <td className="table-cell">{student.section}</td>
                      <td className="table-cell">{student.roll_number}</td>
                      <td className="table-cell">
                        <p className="text-sm">{student.mobile}</p>
                      </td>
                      <td className="table-cell">
                        <span className={`badge ${statusColors[student.status]}`}>
                          {student.status}
                        </span>
                      </td>
                      <td className="table-cell">
                        <button 
                          className="btn btn-ghost btn-sm"
                          onClick={(e) => {
                            e.stopPropagation()
                            setSelectedStudent(student)
                          }}
                        >
                          <MoreHorizontal className="h-4 w-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            <div className="flex items-center justify-between px-6 py-4 border-t">
              <p className="text-sm text-muted-foreground">
                Showing 1 to {filteredStudents.length} of {mockStudents.length} students
              </p>
              <div className="flex items-center gap-2">
                <button 
                  className="btn btn-secondary btn-sm"
                  disabled={page === 1}
                >
                  <ChevronLeft className="mr-1 h-4 w-4" />
                  Previous
                </button>
                <button className="btn btn-secondary btn-sm">
                  Next
                  <ChevronRight className="ml-1 h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Student Details Sidebar */}
        <div className="lg:col-span-1">
          {selectedStudent ? (
            <div className="card sticky top-6">
              <div className="text-center pb-6 border-b">
                <div className="mx-auto mb-4 h-20 w-20 rounded-full bg-primary flex items-center justify-center">
                  <span className="text-2xl font-bold text-primary-foreground">
                    {selectedStudent.first_name[0]}{selectedStudent.last_name[0]}
                  </span>
                </div>
                <h3 className="text-lg font-semibold">
                  {selectedStudent.first_name} {selectedStudent.last_name}
                </h3>
                <p className="text-sm text-muted-foreground">{selectedStudent.student_id}</p>
                <div className="flex justify-center mt-2">
                  <span className={`badge ${statusColors[selectedStudent.status]}`}>
                    {selectedStudent.status}
                  </span>
                </div>
              </div>

              <div className="py-4 space-y-4">
                <div className="flex items-center gap-3 text-sm">
                  <Users className="h-4 w-4 text-muted-foreground" />
                  <div>
                    <p className="text-muted-foreground">Class</p>
                    <p className="font-medium">{selectedStudent.section}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3 text-sm">
                  <Clock className="h-4 w-4 text-muted-foreground" />
                  <div>
                    <p className="text-muted-foreground">Roll Number</p>
                    <p className="font-medium">{selectedStudent.roll_number}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3 text-sm">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <div>
                    <p className="text-muted-foreground">Mobile</p>
                    <p className="font-medium">{selectedStudent.mobile}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3 text-sm">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <div>
                    <p className="text-muted-foreground">Email</p>
                    <p className="font-medium">{selectedStudent.email}</p>
                  </div>
                </div>
              </div>

              <div className="pt-4 border-t flex gap-2">
                <button className="btn btn-primary flex-1">
                  View Profile
                </button>
                <button className="btn btn-secondary">
                  <MoreHorizontal className="h-4 w-4" />
                </button>
              </div>
            </div>
          ) : (
            <div className="card text-center py-12">
              <Users className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
              <p className="text-muted-foreground">Select a student to view details</p>
            </div>
          )}

          {/* Quick Stats */}
          <div className="card mt-4">
            <div className="card-header">
              <h3 className="card-title text-lg">Quick Stats</h3>
            </div>
            <div className="card-content space-y-3">
              <div className="flex items-center justify-between rounded-lg border p-3">
                <span className="text-sm text-muted-foreground">Total Students</span>
                <span className="font-semibold">{mockStudents.length}</span>
              </div>
              <div className="flex items-center justify-between rounded-lg border p-3">
                <span className="text-sm text-muted-foreground">Active</span>
                <span className="font-semibold text-primary">
                  {mockStudents.filter(s => s.status === 'active').length}
                </span>
              </div>
              <div className="flex items-center justify-between rounded-lg border p-3">
                <span className="text-sm text-muted-foreground">Inactive</span>
                <span className="font-semibold text-muted-foreground">
                  {mockStudents.filter(s => s.status === 'inactive').length}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}